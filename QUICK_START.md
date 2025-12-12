# Quick Start Guide

Get up and running with the Peripage Printer API in 5 minutes!

## Prerequisites Check

Before starting, ensure you have:

- ✅ Go 1.22 or higher installed (`go version`)
- ✅ Git installed (for cloning)
- ✅ Docker & Docker Compose (optional, for containerized deployment)

## 🚀 Option 1: Automated Setup (Recommended)

The easiest way to get started:

```bash
cd /Users/erk/Desktop/PrinceM/web-developer/go/sandbox/demo-peri-page-v6
./start.sh
```

This script will:

1. ✅ Check dependencies
2. ✅ Download Go modules
3. ✅ Install Swagger CLI
4. ✅ Generate API documentation
5. ✅ Run tests
6. ✅ Let you choose run mode (Mock or BLE)

Then just follow the prompts!

---

## 🛠️ Option 2: Manual Setup

If you prefer manual control:

### Step 1: Install Dependencies

```bash
cd /Users/erk/Desktop/PrinceM/web-developer/go/sandbox/demo-peri-page-v6
go mod download
```

### Step 2: Install Swagger CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Step 3: Generate Swagger Documentation

```bash
# If swag is in PATH
swag init -g cmd/server/main.go -o internal/adapters/docs

# Or use full path
~/go/bin/swag init -g cmd/server/main.go -o internal/adapters/docs
```

### Step 4: Set Up Environment

```bash
cp .env.example .env
# Edit .env if needed (optional for development)
```

### Step 5: Run the Server

**With Mock Printer (No Hardware):**

```bash
export PRINTER_TYPE=mock
go run cmd/server/main.go
```

**With BLE Printer (Requires Hardware):**

```bash
export PRINTER_TYPE=ble
export PRINTER_DEVICE_NAME=Peripage
go run cmd/server/main.go
```

---

## 🐳 Option 3: Docker Setup

### Development (Mock Printer)

```bash
docker-compose -f docker-compose.dev.yml up --build
```

### Production (BLE Printer - Linux only)

```bash
docker-compose up --build
```

---

## ✅ Verify Installation

### 1. Check Server is Running

Visit in your browser:

```
http://localhost:8080/health
```

You should see:

```json
{
  "status": "healthy"
}
```

### 2. Check Swagger Documentation

Visit in your browser:

```
http://localhost:8080/swagger/index.html
```

You should see the interactive API documentation.

### 3. Test the Print Endpoint

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello, Peripage!"}'
```

Expected response:

```json
{
  "success": true,
  "message": "Print job completed successfully"
}
```

With mock printer, you'll see output in the server logs:

```
[peripage] === MOCK PRINTER OUTPUT ===
Hello, Peripage!
[peripage] === END MOCK PRINTER OUTPUT ===
```

---

## 📚 Next Steps

### Learn More

- Read `README.md` for comprehensive documentation
- Check `API_TESTING.md` for API examples
- Review `PROJECT_STRUCTURE.md` to understand the architecture
- See `BLE_IMPLEMENTATION_TODO.md` for BLE development

### Try the API

**Print JSON data:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "name": "Test User",
      "order": 12345,
      "items": ["A", "B", "C"]
    }
  }'
```

**Use Swagger UI:**

1. Go to http://localhost:8080/swagger/index.html
2. Click "Try it out" on the POST /print endpoint
3. Enter your JSON
4. Click "Execute"

### Development

**Run tests:**

```bash
go test ./...
```

**Build binary:**

```bash
go build -o bin/peripage-server ./cmd/server
./bin/peripage-server
```

**Use Makefile:**

```bash
make help     # See all available commands
make run      # Run with mock printer
make test     # Run tests
make build    # Build binary
```

---

## 🐛 Troubleshooting

### "swag: command not found"

The swag binary is in `~/go/bin/`. Either:

- Add `~/go/bin` to your PATH: `export PATH=$PATH:~/go/bin`
- Use full path: `~/go/bin/swag init -g cmd/server/main.go -o internal/adapters/docs`

### "Port 8080 already in use"

Change the port in `.env`:

```bash
PORT=3000
```

Or set environment variable:

```bash
export PORT=3000
```

### "Cannot find Bluetooth device"

1. Ensure Bluetooth is enabled
2. Check printer is powered on
3. Verify device name matches exactly
4. Try increasing scan timeout in `.env`:
   ```
   BLE_SCAN_TIMEOUT=30s
   ```

### "Docker cannot access Bluetooth"

- On Linux: Use `docker-compose.yml` (production)
- On macOS/Windows: Use `docker-compose.dev.yml` with mock printer
- Docker Desktop doesn't support Bluetooth passthrough

---

## 🎯 Common Use Cases

### Development Without Hardware

```bash
# Use mock printer for development
export PRINTER_TYPE=mock
go run cmd/server/main.go
```

### Testing with Real Printer

```bash
# Connect to actual Peripage device
export PRINTER_TYPE=ble
go run cmd/server/main.go
```

### Production Deployment

```bash
# Build and run in Docker
docker-compose up -d
```

---

## 📞 Need Help?

- Check the documentation files in this project
- Review TODO comments in the code
- Test with mock printer first before BLE
- Use Swagger UI for interactive testing

---

## 🎉 You're Ready!

Your Peripage Printer API is now running. Start making print requests and building your integration!

**Happy Printing! 🖨️**
