#!/bin/bash

# Real-Time Order Processor - Profiling Starter Script
# This script helps you get started with profiling the order processor

echo "🚀 Starting Real-Time Order Processor Profiling Demo"
echo "=================================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if the application is already running
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "⚠️  Application is already running on port 8080"
    echo "   Stopping existing instance..."
    pkill -f "go run cmd/main.go" || true
    sleep 2
fi

# Create profiles directory
mkdir -p profiles

echo "📁 Created profiles directory"

# Start the application in background
echo "🏃 Starting application..."
go run cmd/main.go &
APP_PID=$!

# Wait for application to start
echo "⏳ Waiting for application to start..."
for i in {1..10}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Application started successfully"
        break
    fi
    sleep 1
done

# Check if application started
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Failed to start application"
    kill $APP_PID 2>/dev/null
    exit 1
fi

# Start load generation in background
echo "📊 Starting load generation..."
go run tools/load_test.go &
LOAD_PID=$!

# Wait a bit for load to build up
echo "⏳ Waiting for load to build up..."
sleep 5

# Function to capture profiles
capture_profiles() {
    local scenario=$1
    local duration=$2
    
    echo "📸 Capturing profiles for $scenario scenario..."
    
    # CPU Profile
    echo "   🔥 Capturing CPU profile..."
    curl -s "http://localhost:8080/profile/cpu?duration=${duration}" -o "profiles/cpu_${scenario}.prof"
    
    # Memory Profile
    echo "   🧠 Capturing memory profile..."
    curl -s http://localhost:8080/profile/memory -o "profiles/memory_${scenario}.prof"
    
    # Goroutine Profile
    echo "   🔄 Capturing goroutine profile..."
    curl -s http://localhost:8080/profile/goroutines -o "profiles/goroutine_${scenario}.prof"
    
    # Block Profile
    echo "   🚧 Capturing block profile..."
    curl -s http://localhost:8080/profile/block -o "profiles/block_${scenario}.prof"
    
    # Mutex Profile
    echo "   🔒 Capturing mutex profile..."
    curl -s http://localhost:8080/profile/mutex -o "profiles/mutex_${scenario}.prof"
    
    echo "   ✅ Profiles captured for $scenario"
}

# Capture profiles for different scenarios
capture_profiles "normal" "10s"
sleep 5

capture_profiles "high_load" "15s"
sleep 5

capture_profiles "burst" "5s"

# Stop load generation
echo "🛑 Stopping load generation..."
kill $LOAD_PID 2>/dev/null

# Show available profiles
echo ""
echo "📋 Available profiles:"
ls -la profiles/

echo ""
echo "🔍 To analyze profiles, use:"
echo "   go tool pprof profiles/cpu_normal.prof"
echo "   go tool pprof -http=:8081 profiles/cpu_normal.prof"
echo ""
echo "📖 For detailed instructions, see GO_PROFILING_GUIDE.md"
echo ""
echo "🌐 Profiling web interface: http://localhost:8080/debug/pprof/"
echo "🏥 Health check: http://localhost:8080/health"
echo "📊 Statistics: http://localhost:8080/stats"
echo ""

# Ask if user wants to keep application running
read -p "🤔 Keep application running? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "🛑 Stopping application..."
    kill $APP_PID 2>/dev/null
    echo "✅ Application stopped"
else
    echo "✅ Application running on port 8080"
    echo "   Use Ctrl+C to stop when done"
fi

echo ""
echo "🎉 Profiling demo completed!"
echo "   Check the profiles/ directory for captured data"
