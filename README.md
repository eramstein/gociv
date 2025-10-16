# gociv

Minimal Go + Raylib render loop for a strategy game prototype.

## Prerequisites

- Go 1.22+
- Windows with `raylib.dll` in the project root (already provided)
- A C toolchain for cgo (e.g., MinGW-w64 via MSYS2 or LLVM/clang)

## Run

```bash
# from project root
go run .
```

## Build

```bash
go build -o gociv.exe .
```

If you get cgo or linker errors, ensure your C toolchain is installed and on PATH.
