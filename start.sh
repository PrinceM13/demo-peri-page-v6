#!/bin/bash

# Quick Start Script for Peripage Printer API

set -e

echo "🚀 Peripage Printer API - Quick Start"
echo "======================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.22 or higher."
    exit 1
fi

echo "✅ Go is installed: $(go version)"
echo ""

# Check if .env file exists, if not create from example
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it if needed."
else
    echo "✅ .env file already exists"
fi
echo ""

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
echo "✅ Dependencies downloaded"
echo ""

# Install swag if not present
if ! command -v swag &> /dev/null && ! [ -f ~/go/bin/swag ]; then
    echo "📦 Installing swag for Swagger documentation..."
    go install github.com/swaggo/swag/cmd/swag@latest
    echo "✅ Swag installed"
else
    echo "✅ Swag already installed"
fi
echo ""

# Generate Swagger docs
echo "📚 Generating Swagger documentation..."
if command -v swag &> /dev/null; then
    swag init -g cmd/server/main.go -o internal/adapters/docs
elif [ -f ~/go/bin/swag ]; then
    ~/go/bin/swag init -g cmd/server/main.go -o internal/adapters/docs
fi
echo "✅ Swagger documentation generated"
echo ""

# Run tests
echo "🧪 Running tests..."
go test ./...
echo "✅ Tests passed"
echo ""

# Ask user which mode to run
echo "Select run mode:"
echo "1) Mock Printer (Development - no hardware needed)"
echo "2) BLE Printer (Production - requires Peripage hardware)"
echo "3) Docker Compose (Development)"
echo "4) Docker Compose (Production)"
echo "5) Exit"
echo ""
read -p "Enter choice (1-5): " choice

case $choice in
    1)
        echo ""
        echo "🖨️  Starting with Mock Printer..."
        echo "📍 Server will be available at: http://localhost:8080"
        echo "📚 Swagger docs: http://localhost:8080/swagger/index.html"
        echo ""
        export PRINTER_TYPE=mock
        go run cmd/server/main.go
        ;;
    2)
        echo ""
        echo "🖨️  Starting with BLE Printer..."
        echo "⚠️  Make sure your Peripage printer is powered on and in range"
        echo "📍 Server will be available at: http://localhost:8080"
        echo "📚 Swagger docs: http://localhost:8080/swagger/index.html"
        echo ""
        export PRINTER_TYPE=ble
        go run cmd/server/main.go
        ;;
    3)
        echo ""
        echo "🐳 Starting with Docker Compose (Development)..."
        docker-compose -f docker-compose.dev.yml up --build
        ;;
    4)
        echo ""
        echo "🐳 Starting with Docker Compose (Production)..."
        echo "⚠️  Requires Linux with Bluetooth access"
        docker-compose up --build
        ;;
    5)
        echo "👋 Goodbye!"
        exit 0
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac
