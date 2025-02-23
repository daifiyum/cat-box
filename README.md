## cat-box

A Windows tray app for sing-box.

Use：

```
# Build from source
go build -ldflags="-w -s -H=windowsgui" -o cat-box.exe ./

# Add sing-box.exe
resources/core/sing-box.exe

# Double click exe to run
```

Run Parameters:

`-workspace`: Enable when running with an absolute path, e.g., `E:\cat-box\cat-box.exe -workspace=true`.

`-port`: Customize the backend service port, e.g., `.\cat-box.exe -port=3001`.

Note: For my personal use only, the front-end subscription management is not open-sourced, and no executable program is provided.

