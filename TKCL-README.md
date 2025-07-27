
# 🛠 Installing Tcl/Tk (tkcl) for GUI Support

This guide explains how to install Tcl/Tk on various operating systems so that `modernc.org/tk` (used in this project) can function correctly.

---

## 🪟 Windows 10 / 11

1. **Download ActiveTcl** (recommended):
   - Visit: [https://www.activestate.com/products/tcl/downloads/](https://www.activestate.com/products/tcl/downloads/)
   - Download and install the latest **ActiveTcl Community Edition**.

2. **Verify Installation**:
   ```powershell
   tclsh
   ```

3. **Add to PATH** (if not already):
   - Ensure `tclsh.exe` and `wish.exe` are in your system PATH (e.g., `C:\Tcl\bin`).

---

## 🍎 macOS

Tcl/Tk is preinstalled on macOS, but updating is recommended.

1. **Check if Tcl is available**:
   ```bash
   tclsh
   ```

2. **Optional: Install latest Tcl/Tk via Homebrew**:
   ```bash
   brew install tcl-tk
   ```

3. **Configure environment (if needed)**:
   ```bash
   export PATH="/opt/homebrew/opt/tcl-tk/bin:$PATH"
   export LDFLAGS="-L/opt/homebrew/opt/tcl-tk/lib"
   export CPPFLAGS="-I/opt/homebrew/opt/tcl-tk/include"
   ```

---

## 🐧 Linux Distributions

### Ubuntu / Debian / Linux Mint

```bash
sudo apt update
sudo apt install tk tcl
```

### Fedora

```bash
sudo dnf install tcl tk
```

### Arch Linux / Manjaro

```bash
sudo pacman -S tk tcl
```

### OpenSUSE

```bash
sudo zypper install tcl tk
```

---

## ✅ Verify Installation

Run:
```bash
tclsh
```

If you get a Tcl prompt (`%`), Tcl/Tk is installed correctly.

---

## 💡 Notes

- Ensure `tclsh` and `wish` are in your system PATH.
- `modernc.org/tk` depends on your system Tcl/Tk — no extra Go bindings needed.
- You only need to install Tcl/Tk once per system.
