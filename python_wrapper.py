import platform
import subprocess
import sys

def get_binary_path():
    if platform.system().lower() == "windows":
        return "./tls-trouble.exe"
    return "./tls-trouble_linux"

def fetch(url):
    binary = get_binary_path()
    try:
        result = subprocess.run([binary, url], capture_output=True, text=True, check=True)
        return result.stdout
    except subprocess.CalledProcessError as e:
        return f"[-] error during fetch: {e.stderr}"
    except FileNotFoundError:
        return f"[-] error: binary not found at {binary}."

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"usage: python {sys.argv[0]} <url>")
        sys.exit(1)

    target_url = sys.argv[1]    
    output = fetch(target_url)
    print(output)