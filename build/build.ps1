# build.ps1
# PowerShell 脚本用于编译资源文件和 Go 程序

# 编译资源文件
Write-Output "Compiling resource file..."
# 使用 windres 编译 cat-box.rc 文件到 cat-box.syso
windres -o cat-box.syso cat-box.rc

# 检查编译结果，如果失败则退出脚本
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to compile resource file."
    exit $LASTEXITCODE
}

# 复制到当前目录下的 cat-box.syso
Copy-Item -Path .\cat-box.syso -Destination ..\cat-box.syso

# 编译 Go 程序
Write-Output "Building Go program..."
# 使用 go build 构建 Go 程序，指定 ldflags 为 -w -s -H=windowsgui，输出为 cat-box.exe
go build -ldflags="-w -s -H=windowsgui" -o cat-box.exe ..

# 检查构建结果，如果失败则退出脚本
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build Go program."
    exit $LASTEXITCODE
}

# 清理中间文件
Remove-Item ..\cat-box.syso

# 复制资源文件到输出目录
Copy-Item -Path ..\frontend\dist\* -Destination ..\resources\ui\sub -Recurse -Force
Copy-Item -Path ..\resources -Destination .\ -Recurse -Force

Write-Output "Build succeeded. Output: cat-box.exe"
