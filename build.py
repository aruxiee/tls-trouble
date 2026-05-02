import os
import platform
import subprocess
import sys

def build_tool():
    current_os = platform.system().lower()
    
    if current_os == "windows":
        executable_name = "tls-trouble.exe"
    elif current_os == "linux":
        executable_name = "tls-trouble_linux"
    else:
        executable_name = f"tls-trouble_{current_os}"

    print(f"[*] detected os: {current_os}")

    try:
        print("[*] checking and installing dependencies...")
        subprocess.run(["go", "mod", "tidy"], check=True)
        
        print(f"[*] building {executable_name}...")
        subprocess.run(["go", "build", "-o", executable_name, "main.go"], check=True)
        
        print(f"[+] generated: {executable_name}")
        
    except subprocess.CalledProcessError as e:
        print(f"[-] error occurred during {e.cmd[1]}: build failed.")
        sys.exit(1)
    except FileNotFoundError:
        print("[-] error: 'go' command not found.")
        sys.exit(1)

if __name__ == "__main__":
    build_tool()