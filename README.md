# Peripage A6 Thermal Printer API

A modular Go microservice using Hexagonal Architecture to print text and JSON data to a Peripage A6 mini thermal printer over Bluetooth LE.

## 🏗️ Architecture

This project follows **Hexagonal Architecture (Ports & Adapters)** principles:

```
┌─────────────────────────────────────────────────┐
│              API Layer (Gin)                    │
│          /internal/adapters/api/                │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│           Core Domain Layer                     │
│         /internal/core/                         │
│  • Printer Interface (Port)                     │
│  • PrintService (Business Logic)                │
└────────────────┬────────────────────────────────┘
                 │
      ┌──────────┴──────────┐
      │                     │
┌─────▼──────┐      ┌──────▼─────┐
│ BLEPrinter │      │MockPrinter │
│  Adapter   │      │  Adapter   │
└────────────┘      └────────────┘
```

## 📁 Project Structure

```
demo-peri-page-v6/
├── cmd/
│   └── server/
│       └── main.go              # Entry point with DI
├── internal/
│   ├── core/
│   │   ├── ports.go             # Printer interface
│   │   └── service.go           # PrintService
│   ├── adapters/
│   │   ├── api/
│   │   │   ├── handler.go       # Gin HTTP handlers
│   │   │   └── router.go        # Route configuration
│   │   ├── docs/
│   │   │   └── docs.go          # Swagger docs (generated)
│   │   └── printer/
│   │       ├── ble.go           # BLE printer adapter
│   │       └── mock.go          # Mock printer adapter
│   └── config/
│       └── config.go            # Configuration loader
├── Dockerfile
├── docker-compose.yml           # Production (with BLE)
├── docker-compose.dev.yml       # Development (mock)
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker & Docker Compose (for containerized deployment)
- Bluetooth LE support (for hardware printer)
- Peripage A6 thermal printer (optional, can use mock)

### Installation

1. **Clone the repository:**

   ```bash
   cd /Users/erk/Desktop/PrinceM/web-developer/go/sandbox/demo-peri-page-v6
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Install Swagger CLI (for generating docs):**

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

4. **Generate Swagger documentation:**

   ```bash
   swag init -g cmd/server/main.go -o internal/adapters/docs
   ```

5. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env as needed
   ```

### Configuration

Edit `.env` file:

```bash
# Server Configuration
PORT=8080

# Printer Configuration
PRINTER_TYPE=mock          # Options: mock, ble
PRINTER_DEVICE_NAME=Peripage
PRINTER_TIMEOUT=30s

# BLE Configuration
BLE_SCAN_TIMEOUT=10s
```

## 🔧 Running the Application

### Local Development (Mock Printer)

```bash
# Using mock printer (no hardware needed)
export PRINTER_TYPE=mock
go run cmd/server/main.go
```

### With Real Hardware (BLE Printer)

```bash
# Using BLE printer (requires hardware)
export PRINTER_TYPE=ble
export PRINTER_DEVICE_NAME=Peripage
go run cmd/server/main.go
```

### Using Docker Compose

**Development (Mock Printer):**

```bash
docker-compose -f docker-compose.dev.yml up --build
```

**Production (BLE Printer):**

```bash
docker-compose up --build
```

## 📡 API Endpoints

### Print Text or JSON

**Endpoint:** `POST /print`

**Request Body:**

```json
{
  "text": "Hello, World!",
  "data": {
    "optional": "json",
    "key": "value"
  }
}
```

**Behavior:**

- If `data` is provided, it will be pretty-printed as JSON and sent to printer
- If only `text` is provided, plain text will be printed
- At least one field must be present

**Response (Success):**

```json
{
  "success": true,
  "message": "Print job completed successfully"
}
```

**Response (Error):**

```json
{
  "error": "Print failed: connection timeout"
}
```

### Health Check

**Endpoint:** `GET /health`

**Response:**

```json
{
  "status": "healthy"
}
```

### Swagger Documentation

Access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🧪 Testing the API

### Using cURL

**Print plain text:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello from Peripage!"}'
```

**Print JSON data:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "name": "John Doe",
      "email": "john@example.com",
      "items": ["item1", "item2"]
    }
  }'
```

**Health check:**

```bash
curl http://localhost:8080/health
```

## 🧪 Testing

### Run All Tests
```bash
make test
# or
go test ./... -v
```

### Run Tests with Coverage
```bash
make test-cover
# Generates coverage.html report
```

### Test Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| **internal/core** | 90.0% | ✅ Excellent |
| **internal/adapters/api** | 72.7% | ✅ Good |
| **internal/adapters/printer** | 10.3% | ⚠️ Low (BLE pending) |

### Test Stack

- ✅ **testing** - Go's built-in framework
- ✅ **testify** - Assertions and mocks
- ✅ **httptest** - HTTP handler testing
- ✅ **Build tags** - Separate unit/integration tests

**Full testing guide:** See [TESTING.md](TESTING.md)

**Quick test commands:**
```bash
make test              # Run all tests
make test-cover        # With coverage report
make test-race         # With race detector
make test-integration  # Integration tests (requires hardware)
```

## 🔌 Printer Adapters

### Mock Printer

- **Purpose:** Development and testing without hardware
- **Behavior:** Prints output to stdout
- **Usage:** Set `PRINTER_TYPE=mock`

### BLE Printer

- **Purpose:** Production use with actual Peripage A6 hardware
- **Requirements:** Bluetooth LE support, paired Peripage device
- **Usage:** Set `PRINTER_TYPE=ble`

**⚠️ Important TODOs for BLE Implementation:**

The BLE adapter contains placeholder logic. To complete the implementation:

1. **Device Discovery** (`ble.go:Connect`)

   - Verify exact device name format for Peripage A6
   - Test scanning and device detection

2. **Handshake Protocol** (`ble.go:performHandshake`)

   - Discover Peripage services and characteristics UUIDs
   - Implement initialization commands
   - Enable notifications if required

3. **Text-to-Bitmap Rendering** (`ble.go:textToBitmap`)

   - Implement font rendering (384px width for A6)
   - Convert to 1-bit monochrome bitmap
   - Format according to Peripage protocol

4. **Data Transmission** (`ble.go:sendBitmap`)
   - Implement packet protocol (typically 512-byte max)
   - Add headers/checksums if required
   - Send print trigger command

## 🐳 Docker Deployment

### Development Mode (Mock Printer)

```bash
docker-compose -f docker-compose.dev.yml up -d
```

- Uses regular networking
- No hardware access required
- Perfect for Docker Desktop on macOS/Windows

### Production Mode (BLE Printer)

```bash
docker-compose up -d
```

- Uses `network_mode: host` for Bluetooth access
- Requires `privileged: true`
- Mounts `/var/run/dbus` and `/dev` for hardware access
- Best on Linux systems with direct Bluetooth access

### Switch Between Adapters

Change the `PRINTER_TYPE` environment variable in `docker-compose.yml`:

```yaml
environment:
  - PRINTER_TYPE=mock # or 'ble'
```

## 🔍 Project Features

✅ **Hexagonal Architecture** - Clean separation of concerns  
✅ **Dependency Injection** - Configured in `main.go`  
✅ **Multiple Adapters** - Mock and BLE printer implementations  
✅ **RESTful API** - Using Gin framework  
✅ **Swagger Documentation** - Auto-generated API docs  
✅ **Docker Support** - Full containerization  
✅ **Environment Configuration** - Flexible config via env vars  
✅ **Graceful Shutdown** - Proper cleanup on exit  
✅ **Health Checks** - Docker health monitoring

## 🛠️ Development

### Generate Swagger Docs

After modifying API handlers:

```bash
swag init -g cmd/server/main.go -o internal/adapters/docs
```

### Build Binary

```bash
go build -o bin/peripage-server ./cmd/server
```

### Run Tests

```bash
go test ./...
```

### Code Structure Guidelines

- **Core Domain** (`internal/core/`): Business logic, no external dependencies
- **Ports** (`internal/core/ports.go`): Interfaces defining boundaries
- **Adapters** (`internal/adapters/`): Implementations of ports
- **API Layer**: Only depends on core domain
- **Main**: Wires everything together with dependency injection

## 📝 Adding New Features

### Adding a New Printer Adapter

1. Create file in `internal/adapters/printer/`
2. Implement the `core.Printer` interface
3. Add configuration in `internal/config/config.go`
4. Wire it up in `cmd/server/main.go`

### Adding a New API Endpoint

1. Add handler method in `internal/adapters/api/handler.go`
2. Add Swagger annotations
3. Register route in `internal/adapters/api/router.go`
4. Regenerate Swagger docs

## 🐛 Troubleshooting

### Bluetooth Issues

**Problem:** Cannot find Bluetooth device

- Ensure Bluetooth is enabled on host
- Check device name matches exactly
- Increase `BLE_SCAN_TIMEOUT`

**Problem:** Permission denied accessing Bluetooth

- Ensure container runs with `privileged: true`
- Check `/var/run/dbus` is mounted
- Verify user has Bluetooth permissions

### Docker Issues

**Problem:** Container cannot access Bluetooth

- Use `network_mode: host` on Linux
- On macOS/Docker Desktop, use mock printer instead

## 📄 License

MIT

## 👤 Author

PrinceM

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Swaggo](https://github.com/swaggo/swag)
- [TinyGo Bluetooth](https://github.com/tinygo-org/bluetooth)
