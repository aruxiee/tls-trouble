package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/andybalholm/brotli"
	utls "github.com/refraction-networking/utls"
)

type Profile struct {
	Name      string
	HelloID   utls.ClientHelloID
	UserAgent string
}

var BrowserProfiles = []Profile{
	{
		Name:      "Chrome_Windows",
		HelloID:   utls.HelloChrome_120,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	},
	{
		Name:      "Firefox_Windows",
		HelloID:   utls.HelloFirefox_105,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0",
	},
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("usage: %s <url>\n", os.Args[0])
		return
	}

	u, err := url.Parse(args[0])
	if err != nil || u.Host == "" {
		u, _ = url.Parse("https://" + args[0])
	}

	targetHost := u.Hostname()
	targetPath := u.RequestURI()
	targetAddr := targetHost + ":443"

	rand.Seed(time.Now().UnixNano())
	activeProfile := BrowserProfiles[rand.Intn(len(BrowserProfiles))]
	fmt.Fprintf(os.Stderr, "[+] identity: %s\n", activeProfile.Name)

	conn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		return
	}

	config := &utls.Config{ServerName: targetHost}
	tlsConn := utls.UClient(conn, config, utls.HelloCustom)
	
	spec, err := utls.UTLSIdToSpec(activeProfile.HelloID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "spec creation error: %v\n", err)
		return
	}
	
	for i, ext := range spec.Extensions {
		if alpnExt, ok := ext.(*utls.ALPNExtension); ok {
			alpnExt.AlpnProtocols = []string{"http/1.1"}
			spec.Extensions[i] = alpnExt
		}
	}

	if err := tlsConn.ApplyPreset(&spec); err != nil {
		fmt.Fprintf(os.Stderr, "preset error: %v\n", err)
		return
	}

	if err := tlsConn.Handshake(); err != nil {
		fmt.Fprintf(os.Stderr, "handshake error: %v\n", err)
		return
	}

	reqStr := fmt.Sprintf("GET %s HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"User-Agent: %s\r\n"+
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"+
		"Accept-Encoding: gzip, deflate, br\r\n"+
		"Connection: close\r\n\r\n",
		targetPath, targetHost, activeProfile.UserAgent)

	tlsConn.Write([]byte(reqStr))

	reader := bufio.NewReader(tlsConn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing response: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var bodyReader io.ReadCloser
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "br":
		bodyReader = io.NopCloser(brotli.NewReader(resp.Body))
	case "gzip":
		bodyReader, _ = gzip.NewReader(resp.Body)
	default:
		bodyReader = resp.Body
	}
	defer bodyReader.Close()

	_, err = io.Copy(os.Stdout, bodyReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "output stream error: %v\n", err)
	}
}