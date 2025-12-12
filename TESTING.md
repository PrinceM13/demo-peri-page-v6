# Testing Guide

Comprehensive testing documentation for the Peripage Printer API.

## Test Stack (2025 Best Practices)

- **testing** - Go's built-in testing package
- **testify/assert** - Fluent assertions
- **testify/require** - Assertions that stop test on failure
- **testify/mock** - Mock generation and verification
- **httptest** - HTTP handler testing
- **Build tags** - Separate unit tests from integration tests

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Run Tests Verbosely
```bash
go test ./... -v
```

### Run Specific Package Tests
```bash
go test ./internal/core -v
go test ./internal/adapters/api -v
go test ./internal/adapters/printer -v
```

### Run Specific Test
```bash
go test ./internal/core -run TestPrintService_PrintText -v
```

### Run with Coverage Report (HTML)
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # macOS
```

### Run Integration Tests (Requires Hardware)
```bash
go test -tags=integration ./test/integration/... -v
```

### Skip Integration Tests (Default)
```bash
go test ./...  # Automatically skips integration tests
```

### Run Tests in Short Mode
```bash
go test ./... -short
```

## Test Coverage

Current coverage by package:

| Package | Coverage | Status |
|---------|----------|--------|
| **internal/core** | 90.0% | ✅ Excellent |
| **internal/adapters/api** | 72.7% | ✅ Good |
| **internal/adapters/printer** | 10.3% | ⚠️ Low (BLE pending implementation) |

## Test Structure

### Core Domain Tests (`internal/core/service_test.go`)

Tests the business logic layer with mocked dependencies.

**Test Cases:**
- ✅ Service creation
- ✅ Print text (valid, empty, errors)
- ✅ Print JSON (valid, nested, nil, errors)
- ✅ JSON formatting verification
- ✅ Error propagation from printer

**Example:**
```go
func TestPrintService_PrintText(t *testing.T) {
    // Arrange
    mockPrinter := new(mocks.MockPrinter)
    mockPrinter.On("PrintText", "Hello").Return(nil).Once()
    service := NewPrintService(mockPrinter)
    
    // Act
    err := service.PrintText("Hello")
    
    // Assert
    require.NoError(t, err)
    mockPrinter.AssertExpectations(t)
}
```

### API Handler Tests (`internal/adapters/api/handler_test.go`)

Tests HTTP endpoints using `httptest` with mocked service.

**Test Cases:**
- ✅ POST /print with text
- ✅ POST /print with JSON data
- ✅ POST /print with both (data takes precedence)
- ✅ Invalid JSON body
- ✅ Empty request body
- ✅ Printer error propagation
- ✅ GET /health endpoint

**Example:**
```go
func TestHandler_Print(t *testing.T) {
    // Arrange
    mockPrinter := new(mocks.MockPrinter)
    mockPrinter.On("PrintText", "Hello").Return(nil)
    service := &mockPrintService{printer: mockPrinter}
    handler := &Handler{service: service}
    router := setupTestRouter(handler)
    
    body := `{"text":"Hello"}`
    req := httptest.NewRequest("POST", "/print", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Act
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockPrinter.AssertExpectations(t)
}
```

### Adapter Tests

#### Mock Printer Tests (`internal/adapters/printer/mock_test.go`)

Tests the mock printer implementation (always succeeds).

**Test Cases:**
- ✅ Printer creation with/without logger
- ✅ Print various text types
- ✅ Always succeeds (no errors)
- ✅ Log output verification

#### BLE Printer Tests (`internal/adapters/printer/ble_test.go`)

Placeholder tests for BLE printer (requires hardware).

**Build Tag:** `//go:build !integration`

**Test Cases:**
- ⏭️ All tests skipped (hardware required)
- 📝 TODO comments for implementation
- 🔗 Redirects to integration tests

### Integration Tests (`test/integration/ble_integration_test.go`)

Real hardware tests with actual Peripage A6 device.

**Build Tag:** `//go:build integration`

**Test Cases (TODO):**
- ⏭️ Device discovery
- ⏭️ Connection establishment
- ⏭️ Print text
- ⏭️ Print JSON
- ⏭️ Multiple jobs
- ⏭️ Error recovery
- ⏭️ Disconnect

**Run Integration Tests:**
```bash
# Set device name (optional)
export PERIPAGE_DEVICE_NAME=Peripage

# Run tests
go test -tags=integration ./test/integration/... -v

# Run with timeout
go test -tags=integration ./test/integration/... -v -timeout=5m
```

## Test Patterns Used

### Table-Driven Tests
```go
tests := []struct {
    name          string
    input         string
    expectedError string
}{
    {
        name:          "valid input",
        input:         "hello",
        expectedError: "",
    },
    {
        name:          "empty input",
        input:         "",
        expectedError: "cannot be empty",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### Mock Setup Pattern
```go
mockPrinter := new(mocks.MockPrinter)
mockPrinter.On("PrintText", "expected text").Return(nil).Once()
// ... run test ...
mockPrinter.AssertExpectations(t)
```

### HTTP Testing Pattern
```go
router := setupTestRouter(handler)
req := httptest.NewRequest("POST", "/print", body)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: go test ./... -cover
      
      - name: Run tests with race detector
        run: go test ./... -race
```

### Makefile Integration
```makefile
test:
	go test -v ./...

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-race:
	go test ./... -race

test-integration:
	go test -tags=integration ./test/integration/... -v
```

## Mocking Strategy

### Using testify/mock

**Create Mock:**
```go
type MockPrinter struct {
    mock.Mock
}

func (m *MockPrinter) PrintText(text string) error {
    args := m.Called(text)
    return args.Error(0)
}
```

**Setup Expectations:**
```go
mockPrinter := new(MockPrinter)
mockPrinter.On("PrintText", "expected").Return(nil).Once()
mockPrinter.On("PrintText", "error").Return(errors.New("failed")).Once()
```

**Verify Calls:**
```go
mockPrinter.AssertExpectations(t)
mockPrinter.AssertCalled(t, "PrintText", "expected")
mockPrinter.AssertNumberOfCalls(t, "PrintText", 2)
```

## Best Practices

### ✅ DO

- **Use table-driven tests** for multiple scenarios
- **Mock external dependencies** (printers, network, etc.)
- **Test both success and error paths**
- **Use require for critical assertions** (stops test on failure)
- **Use assert for non-critical assertions** (continues test)
- **Name tests descriptively** (`TestServiceName_MethodName_Scenario`)
- **Use build tags** to separate unit and integration tests
- **Keep tests deterministic** (no random values, no time dependencies)
- **Test interfaces, not implementations**

### ❌ DON'T

- **Don't test hardware in unit tests**
- **Don't use real BLE devices in CI**
- **Don't use sleep() for timing** (use mocks/timeouts)
- **Don't test external APIs directly**
- **Don't commit failing tests**
- **Don't skip writing tests** (even placeholders)

## Writing New Tests

### For Core Domain
1. Create mock for dependencies
2. Set up expectations
3. Call method under test
4. Assert results
5. Verify mock expectations

### For API Handlers
1. Create test router with handler
2. Create HTTP request with httptest
3. Create response recorder
4. Execute request
5. Assert response status and body

### For Adapters
1. Test mock adapter fully (no hardware)
2. Create placeholder tests for hardware adapters
3. Use build tags for integration tests
4. Skip tests that need hardware
5. Document TODOs for future implementation

## Debugging Tests

### Run Single Test with Verbose Output
```bash
go test -v -run TestPrintService_PrintText ./internal/core
```

### Print Debug Info
```go
t.Logf("Debug info: %v", value)
t.Log("Step 1 completed")
```

### Fail Test Manually
```go
t.Fatal("Critical error occurred")
t.Error("Non-critical error")
```

### Skip Test Conditionally
```go
if !hardwareAvailable {
    t.Skip("Hardware not available")
}
```

## Test Maintenance

### When to Update Tests

- ✅ When adding new features
- ✅ When fixing bugs
- ✅ When changing APIs
- ✅ When refactoring code
- ✅ Before merging PRs

### Test Checklist

- [ ] All tests pass locally
- [ ] Coverage maintained or improved
- [ ] No flaky tests
- [ ] No skipped tests (except hardware)
- [ ] Descriptive test names
- [ ] Clear assertions
- [ ] Proper cleanup (defers)

## Coverage Goals

| Component | Target | Current |
|-----------|--------|---------|
| Core Domain | >80% | 90.0% ✅ |
| API Layer | >70% | 72.7% ✅ |
| Adapters (Mock) | >80% | ✅ |
| Adapters (BLE) | >50% | ⏳ Pending |
| Overall | >70% | ✅ |

## Future Improvements

- [ ] Add benchmark tests for performance
- [ ] Add fuzz tests for input validation
- [ ] Implement BLE adapter integration tests
- [ ] Add end-to-end tests
- [ ] Set up automated coverage reporting
- [ ] Add mutation testing
- [ ] Create test data factories

## Quick Reference

### Run Commands
```bash
go test ./...                          # All tests
go test ./... -cover                   # With coverage
go test ./... -v                       # Verbose
go test ./... -short                   # Skip slow tests
go test -tags=integration ./test/...   # Integration only
go test -run TestName ./package         # Specific test
go test ./... -race                    # Race detection
go test ./... -count=1                 # Disable cache
```

### Assertions
```go
// Require (stops on failure)
require.NoError(t, err)
require.NotNil(t, value)
require.Equal(t, expected, actual)

// Assert (continues on failure)
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.Contains(t, str, substr)
assert.True(t, condition)
```

### Mock Setup
```go
mock.On("Method", arg).Return(result).Once()
mock.On("Method", mock.Anything).Return(result)
mock.AssertExpectations(t)
mock.AssertCalled(t, "Method", arg)
```

---

**Happy Testing! 🧪**
