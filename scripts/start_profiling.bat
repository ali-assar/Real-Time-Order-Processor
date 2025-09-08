@echo off
REM Real-Time Order Processor - Profiling Starter Script for Windows
REM This script helps you get started with profiling the order processor

echo 🚀 Starting Real-Time Order Processor Profiling Demo
echo ==================================================

REM Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go first.
    pause
    exit /b 1
)

REM Check if the application is already running
curl -s http://localhost:8080/health >nul 2>&1
if %errorlevel% equ 0 (
    echo ⚠️  Application is already running on port 8080
    echo    Stopping existing instance...
    taskkill /f /im go.exe >nul 2>&1
    timeout /t 2 >nul
)

REM Create profiles directory
if not exist profiles mkdir profiles
echo 📁 Created profiles directory

REM Start the application in background
echo 🏃 Starting application...
start /b go run cmd/main.go
timeout /t 3 >nul

REM Wait for application to start
echo ⏳ Waiting for application to start...
for /l %%i in (1,1,10) do (
    curl -s http://localhost:8080/health >nul 2>&1
    if %errorlevel% equ 0 (
        echo ✅ Application started successfully
        goto :app_started
    )
    timeout /t 1 >nul
)

echo ❌ Failed to start application
pause
exit /b 1

:app_started
REM Start load generation in background
echo 📊 Starting load generation...
start /b go run tools/load_test.go

REM Wait a bit for load to build up
echo ⏳ Waiting for load to build up...
timeout /t 5 >nul

REM Function to capture profiles
echo 📸 Capturing profiles for normal scenario...
echo    🔥 Capturing CPU profile...
curl -s "http://localhost:8080/profile/cpu?duration=10s" -o "profiles/cpu_normal.prof"

echo    🧠 Capturing memory profile...
curl -s http://localhost:8080/profile/memory -o "profiles/memory_normal.prof"

echo    🔄 Capturing goroutine profile...
curl -s http://localhost:8080/profile/goroutines -o "profiles/goroutine_normal.prof"

echo    🚧 Capturing block profile...
curl -s http://localhost:8080/profile/block -o "profiles/block_normal.prof"

echo    🔒 Capturing mutex profile...
curl -s http://localhost:8080/profile/mutex -o "profiles/mutex_normal.prof"

echo    ✅ Profiles captured for normal

timeout /t 5 >nul

echo 📸 Capturing profiles for high load scenario...
echo    🔥 Capturing CPU profile...
curl -s "http://localhost:8080/profile/cpu?duration=15s" -o "profiles/cpu_high_load.prof"

echo    🧠 Capturing memory profile...
curl -s http://localhost:8080/profile/memory -o "profiles/memory_high_load.prof"

echo    🔄 Capturing goroutine profile...
curl -s http://localhost:8080/profile/goroutines -o "profiles/goroutine_high_load.prof"

echo    🚧 Capturing block profile...
curl -s http://localhost:8080/profile/block -o "profiles/block_high_load.prof"

echo    🔒 Capturing mutex profile...
curl -s http://localhost:8080/profile/mutex -o "profiles/mutex_high_load.prof"

echo    ✅ Profiles captured for high load

timeout /t 5 >nul

echo 📸 Capturing profiles for burst scenario...
echo    🔥 Capturing CPU profile...
curl -s "http://localhost:8080/profile/cpu?duration=5s" -o "profiles/cpu_burst.prof"

echo    🧠 Capturing memory profile...
curl -s http://localhost:8080/profile/memory -o "profiles/memory_burst.prof"

echo    🔄 Capturing goroutine profile...
curl -s http://localhost:8080/profile/goroutines -o "profiles/goroutine_burst.prof"

echo    🚧 Capturing block profile...
curl -s http://localhost:8080/profile/block -o "profiles/block_burst.prof"

echo    🔒 Capturing mutex profile...
curl -s http://localhost:8080/profile/mutex -o "profiles/mutex_burst.prof"

echo    ✅ Profiles captured for burst

REM Stop load generation
echo 🛑 Stopping load generation...
taskkill /f /im go.exe >nul 2>&1

REM Show available profiles
echo.
echo 📋 Available profiles:
dir profiles\

echo.
echo 🔍 To analyze profiles, use:
echo    go tool pprof profiles/cpu_normal.prof
echo    go tool pprof -http=:8081 profiles/cpu_normal.prof
echo.
echo 📖 For detailed instructions, see GO_PROFILING_GUIDE.md
echo.
echo 🌐 Profiling web interface: http://localhost:8080/debug/pprof/
echo 🏥 Health check: http://localhost:8080/health
echo 📊 Statistics: http://localhost:8080/stats
echo.

REM Ask if user wants to keep application running
set /p keep_running="🤔 Keep application running? (y/n): "
if /i "%keep_running%"=="y" (
    echo ✅ Application running on port 8080
    echo    Use Ctrl+C to stop when done
    pause
) else (
    echo 🛑 Stopping application...
    taskkill /f /im go.exe >nul 2>&1
    echo ✅ Application stopped
)

echo.
echo 🎉 Profiling demo completed!
echo    Check the profiles/ directory for captured data
pause
