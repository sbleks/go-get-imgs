@echo off
echo Building Windows executable...
go build -o go-get-imgs.exe main.go
if %ERRORLEVEL% EQU 0 (
    echo Build successful! Created go-get-imgs.exe
    echo.
    echo Usage: go-get-imgs.exe <csv-file> <url-column-index>   
    echo Example: go-get-imgs.exe sample.csv 3
) else (
    echo Build failed!
)
pause 