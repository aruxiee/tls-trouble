
# 🆘 tls-trouble: Protocol Stealth & Fingerprinting Engine

**tls-trouble** is a cross-platform Go script designed to audit web security filters and bypass protocol-based detection. It utilizes **uTLS** to mimic high-fidelity browser fingerprints and strips **ALPN** extensions, allowing modern servers to downgrade to HTTP/1.1, bypasses binary protocol traps (HTTP/2) and JA3/JA3S/JA4 fingerprinting used by CDNs and WAFs.

⚠️ **Please Note:** This project is strictly for **Educational and Authorized Penetration Testing**. I am not responsible for any of the shenanigans you guys pull.

## 🚀 Technical Overview

Modern security solutions monitor the TLS handshake to identify automated scripts. This tool counters this by masking its identity as a legitimate browser and modifying the handshake negotiation.

*   **Fingerprint Spoof:** Perfectly replicates the `ClientHello` of Chrome 120 and Firefox 105.
*   **Protocol Forcing:** Manually edits TLS extensions to remove `h2` (HTTP/2) support, making a readable HTTP/1.1 stream.
*   **Cross-Platform:** Single Go source required for Windows and Linux.
*   **Scalability:** Can be scaled for high-performance lead database management and data extraction.

## 💎 Implementation Details

You will need the latest version of Go to make the installation happen. Either download the installer from [here](https://golang.org/dl/) or use package managers to install it.

- **Windows:** `choco install golang`
- **Debian:** `sudo apt install golang-go`
- **Arch:** `sudo pacman -S go`
- **MacOS:** `brew install go`

### 🛠️ Build the tool (`build.py`)

This script handles cross-compilation and dependency management.

#### 1. Compilation
Run the builder to download required Go dependencies and generate the Go binary for your current environment.
```python
python3 build.py
```

#### 2. Script Explanation
The builder executes `go mod tidy` to verify the environment before using `go build -o` to output the binary.
*   **Windows:** `tls-trouble.exe`
*   **Linux:** `tls-trouble_linux`

### 💅🏻 Python Wrapper

The Python script acts as a high-level wrapper, allowing you to pass URLs directly from the terminal to the underlying engine.

Pass any URL as an argument to the wrapper. It automatically detects your OS and uses the correct binary.
```bash
python3 python_wrapper.py <url>
```

## 噫 Implementation Logic

### ⚙️ Go Engine (`main.go`)

The core of **tls-trouble** operates as a proxy-mimic. It doesn't just change `User-Agent` string but also re-engineers communication stack to ensure consistency across the OSI model.

### 制 Layer 7: Semantic & Protocol Spoofing (App)
At the highest level, the engine ensures the server receives exactly what it expects from a modern browser:
*   **Header Synchronization:** `User-Agent` is paired strictly with the TLS handshake identity. A Firefox handshake with a Chrome User-Agent is a common mishap that this engine prevents.
*   **ALPN Manipulation:** The engine forcibly strips the `h2` (HTTP/2) protocol from the extension. By doing this, the server is forced to fall back to HTTP/1.1. This bypasses complex HTTP/2-specific fingerprinting (like header frame ordering) and allows for easier data extraction via standard streams.
*   **Encoding Fidelity:** Full support for `br` (Brotli) and `gzip`. The engine can process high-compression responses used by modern websites.

### ｪ Layer 4: TLS & Handshake Integrity (Transport)
This tool also counters advanced ghost detection systems. Most scrapers use the default Go `crypto/tls` library, which has a weak signature.

*   **uTLS Fingerprint Spoof:** The engine uses `refraction-networking/utls` library to construct a **ClientHello** packet that is bit-by-bit identical to a real browser. This includes
    *   **Cipher Suite Ordering:** Matches the exact priority list of encryption algorithms used by Chrome or Firefox.
    *   **Extension Padding:** Replicates the grease values and specific extension lengths that modern browsers use to prevent ossification.
*   **JA3/JA3S/JA4:** By mimicking the TLS handshake, the tool produces legitimate **JA3/JA4 hashes**. Security filters (like Cloudflare or Akamai) see a trusted browser hash instead of a `Go-http-client` hash.
*   **TCP/IP Stack Consistency:** Because it runs as a native binary (compiled via `build.py`), the engine utilizes the host's native networking stack. When combined with the TLS mimicking, it ensures that the packet timing and window sizes appear consistent with a desktop user on Linux or Windows.

### ｧｩ Mechanics of the "Identity Swap"
The script follows a strict execution flow to maintain this mask.
- **Selection:** A random `Profile` is chosen (e.g., `Chrome_Windows`).
- **Socket Creation:** A raw TCP connection is opened.
- **Handshake Customization:** Engine generates a `utls.Config` and applies a `UClient` preset.
- **tls-trouble:** Before the handshake is sent, the script intercepts `ALPNExtension` and deletes the `h2` protocol string.
- **Encrypted Tunnel:** The handshake completes, creating a tunnel that the server believes is a standard browser session.

## 📊 Evasion Matrix

| Detection Layer | Standard Scraper | tls-trouble |
| :--- | :--- | :--- |
| **TLS Fingerprinting** | Uses default Go/Python ciphers (easily flagged). | Mimics Chrome/Firefox handshake ciphers and extensions. |
| **Protocol Traps** | Connects via HTTP/2, which is harder to manipulate. | Forces HTTP/1.1 downgrade via ALPN stripping. |
| **User-Agent Matching** | Handshake doesn't match the header. | Handshake and User-Agent are randomized in pairs. |
| **WAF Filtering** | Blocks non-browser JA3/JA4 signatures. | Depicts a perfect JA3/JA4 signature of a modern browser. |

### 🎀 Why use it?

You should use **tls-trouble** regularly if you are scraping sensitive lead lists or business data, as standard libraries are easily blocked by modern CDNs. This tool ensures your connection is indistinguishable from a regular user browsing the internet.

*   **Fingerprint Spoofing:** With high-fidelity `uTLS` profiles, the tool bypasses JA3 and JA4 blocklists that typically flag automated scripts.
*   **Protocol Downgrade:** The ALPN stripping allows to force traffic into legacy inspection streams, which can bypass modern protocol traps.
*   **Signature Masking:** Synchronizing headers with cipher suites prevents detection from signature-based IDS/IPS systems.
*   **Automated Auditing:** The Python wrapper allows for rapid testing and data collection across thousands of endpoints.

### 🕵️ Pentesting & Security Auditing
Beyond data collection, **tls-trouble** also serves as a tool for Red Team and pentesting scenarios.

*   **WAF Evasion Testing:** Verify if a WAF is actually inspecting handshakes or simply allow-listing specific JA3/JA4 browser hashes.
*   **Protocol Downgrade Attacks:** Simulate scenarios where stripping `h2` from the ALPN extension forces a server into a HTTP/1.1 stream.
*   **Stealthy Reconnaissance:** Perform info gathering (like fetching `.env` or `.git` files) without leaving automated scanner traces in logs.
*   **Cipher Suite Auditing:** Tweak profiles to use deprecated ciphers to know if a server's TLS config follows modern hardening standards.

### 🤔 How do I use it?

*   **Profiles:** Add new browsers to the `BrowserProfiles` slice in `main.go` to expand your identity rotation.
*   **Integration:** Use the `fetch()` function in the wrapper to integrate the Go engine into your lead management or pentesting scripts.
*   **Headers:** Tweak the `reqStr` in Go to add custom cookies, referral headers, specific session tokens, etc. for site interaction.

---

<p align="center">
  With ❤️ by <b>Aradhya</b>
</p>
