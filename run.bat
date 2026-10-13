@echo off
REM Builds the menu service and starts it on http://localhost:8080
REM Double-click this file, then open the address in your browser.

cd /d "%~dp0"

echo Building...
go build -o menu.exe .
if errorlevel 1 (
    echo.
    echo Build failed. Check the errors above.
    pause
    exit /b 1
)

echo.
echo ==========================================
echo   Panel : http://localhost:8080/login
echo   Demo  : http://localhost:8080/m/demo
echo.
echo   Press Ctrl+C to stop the server.
echo ==========================================
echo.

menu.exe
pause
