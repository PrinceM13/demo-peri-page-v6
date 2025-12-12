# BLE Implementation TODO List

This document outlines the tasks needed to complete the Bluetooth LE integration with the Peripage A6 thermal printer.

## 🔍 Research Phase

### Device Information

- [ ] Power on Peripage A6 and verify exact Bluetooth device name
- [ ] Document the advertised service UUIDs
- [ ] Check if device requires pairing before connection
- [ ] Test connection range and stability

### Protocol Documentation

- [ ] Find official Peripage A6 protocol documentation
- [ ] Reverse engineer protocol if documentation not available
- [ ] Document all command codes and packet formats
- [ ] Identify acknowledgment/status response format

## 🔧 Implementation Tasks

### 1. Device Discovery (`ble.go:Connect`)

**File:** `internal/adapters/printer/ble.go`  
**Function:** `Connect()`

**Current State:** Placeholder implementation with basic scanning

**Tasks:**

- [ ] Verify device name filter (might be "Peripage", "MX06", or similar)
- [ ] Add MAC address filtering as alternative to name
- [ ] Handle multiple devices with same name
- [ ] Implement retry logic for failed connections
- [ ] Add timeout handling
- [ ] Log RSSI (signal strength) for debugging

**Code Location:** Lines ~35-75

**Testing:**

```bash
# Test scanning
go run cmd/server/main.go  # with PRINTER_TYPE=ble
# Check logs for discovered devices
```

---

### 2. Service & Characteristic Discovery (`ble.go:performHandshake`)

**File:** `internal/adapters/printer/ble.go`  
**Function:** `performHandshake()`

**Current State:** Empty placeholder

**Tasks:**

- [ ] Discover all available services
- [ ] Identify printer service UUID (commonly 0x18F0 or custom)
- [ ] Find write characteristic UUID
- [ ] Find notify/indicate characteristic UUID (if any)
- [ ] Subscribe to notifications for status updates
- [ ] Document all discovered UUIDs in code comments

**Common Printer Service UUIDs to try:**

- `0000FE00-0000-1000-8000-00805F9B34FB` (Generic)
- `000018F0-0000-1000-8000-00805F9B34FB` (Printer Service)
- Custom vendor UUIDs

**Code Location:** Lines ~90-105

**Example Code:**

```go
services, err := b.device.DiscoverServices(nil)
if err != nil {
    return fmt.Errorf("failed to discover services: %w", err)
}

for _, service := range services {
    b.logger.Printf("Service UUID: %s", service.UUID().String())

    chars, err := service.DiscoverCharacteristics(nil)
    if err != nil {
        continue
    }

    for _, char := range chars {
        b.logger.Printf("  Characteristic UUID: %s, Properties: %v",
            char.UUID().String(), char.Properties())
    }
}
```

---

### 3. Initialization Commands (`ble.go:performHandshake`)

**File:** `internal/adapters/printer/ble.go`  
**Function:** `performHandshake()`

**Current State:** Not implemented

**Tasks:**

- [ ] Send printer wake-up command
- [ ] Set print density (darkness level)
- [ ] Set paper type (if applicable)
- [ ] Query printer status
- [ ] Verify printer is ready
- [ ] Handle error responses

**Common Initialization Commands:**

```
Wake-up:     [0x10, 0xFF, 0xFE, 0x01]  // Example
Set Density: [0x1D, 0x7C, 0x50, 0x02]  // Example, 0x02 = level
Status:      [0x1B, 0x76]              // Example
```

**Code Location:** Lines ~90-105

**Testing:**

- Use a BLE sniffer app (nRF Connect) to capture commands from official app
- Send test commands and observe printer response

---

### 4. Text to Bitmap Conversion (`ble.go:textToBitmap`)

**File:** `internal/adapters/printer/ble.go`  
**Function:** `textToBitmap()`

**Current State:** Returns error, not implemented

**Tasks:**

- [ ] Install font rendering library (`golang.org/x/image/font`)
- [ ] Choose font (or embed TrueType font)
- [ ] Create image with correct width (384px for Peripage A6)
- [ ] Render text onto image
- [ ] Calculate height based on text content
- [ ] Convert to 1-bit monochrome bitmap
- [ ] Apply dithering if needed for quality
- [ ] Format as byte array (printer format)
- [ ] Add support for text wrapping
- [ ] Add support for different font sizes

**Peripage A6 Specs:**

- Width: 384 pixels (48 bytes)
- Height: Variable
- Format: 1-bit monochrome
- Byte order: MSB first (usually)

**Code Location:** Lines ~120-140

**Example Implementation:**

```go
import (
    "image"
    "image/color"
    "image/draw"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
)

const printerWidthPx = 384
const printerWidthBytes = 48

func (b *BLEPrinter) textToBitmap(text string) ([]byte, error) {
    // TODO: Implement proper font rendering
    // 1. Create white image
    // 2. Draw black text
    // 3. Convert to 1-bit bitmap
    // 4. Return byte array

    return nil, fmt.Errorf("not implemented")
}
```

**Testing:**

```go
// Test bitmap generation
bitmap, err := printer.textToBitmap("Test")
// Verify bitmap dimensions and format
// Display bitmap as image for visual verification
```

---

### 5. Bitmap Transmission (`ble.go:sendBitmap`)

**File:** `internal/adapters/printer/ble.go`  
**Function:** `sendBitmap()`

**Current State:** Returns error, not implemented

**Tasks:**

- [ ] Determine maximum packet size (typically 512 bytes)
- [ ] Split bitmap into appropriate chunks
- [ ] Add packet headers if required by protocol
- [ ] Calculate and add checksums if required
- [ ] Send packets sequentially
- [ ] Wait for acknowledgments if needed
- [ ] Handle retransmission on failure
- [ ] Send final "print" command to trigger printing
- [ ] Add progress logging

**Packet Structure (Example):**

```
[Header: 0x02] [Length: 2 bytes] [Data: N bytes] [Checksum: 1 byte]
```

**Code Location:** Lines ~143-160

**Example Implementation:**

```go
const maxPacketSize = 512

func (b *BLEPrinter) sendBitmap(bitmap []byte) error {
    // Split into packets
    for offset := 0; offset < len(bitmap); offset += maxPacketSize {
        end := offset + maxPacketSize
        if end > len(bitmap) {
            end = len(bitmap)
        }

        packet := bitmap[offset:end]
        // TODO: Add headers, checksum
        // TODO: Write to characteristic
        // TODO: Wait for ACK
    }

    // TODO: Send print trigger command
    return nil
}
```

**Testing:**

- Send small test bitmap first (e.g., 10 lines)
- Verify each packet reaches printer
- Check for correct printing output

---

## 🧪 Testing Strategy

### Phase 1: Discovery & Connection

1. Run with BLE adapter enabled
2. Verify device is discovered
3. Verify connection succeeds
4. Log all services and characteristics

### Phase 2: Handshake

1. Send initialization commands
2. Verify printer responds
3. Check printer status

### Phase 3: Bitmap Rendering

1. Generate test bitmap offline
2. Save as image file to verify visually
3. Test with different text lengths
4. Test with special characters

### Phase 4: Transmission

1. Send test bitmap (solid black square)
2. Send text bitmap
3. Test with long text (pagination)
4. Test error handling

### Phase 5: Integration

1. Test via API endpoint
2. Test multiple print jobs
3. Test disconnection/reconnection
4. Stress test with rapid requests

---

## 🔨 Tools & Resources

### Development Tools

- **nRF Connect** (iOS/Android): BLE scanning and sniffing
- **Wireshark**: Capture BLE packets (requires hardware)
- **BLE Explorer**: Alternative BLE tool for development

### Go Libraries

- `tinygo.org/x/bluetooth`: Already included
- `golang.org/x/image/font`: For text rendering
- `github.com/golang/freetype`: Advanced font rendering

### Documentation

- [ ] Search for "Peripage A6 protocol" documentation
- [ ] Check Peripage SDK if available
- [ ] Review similar thermal printer protocols (ESC/POS)
- [ ] Study existing open-source implementations

### Reverse Engineering

If official docs not available:

1. Install official Peripage app
2. Use BLE sniffer to capture communication
3. Document packet sequences
4. Identify command patterns
5. Test commands individually

---

## 📋 Testing Checklist

Before marking BLE implementation complete:

- [ ] Device discovery works reliably
- [ ] Connection is stable
- [ ] Can reconnect after disconnection
- [ ] Text renders correctly
- [ ] Print quality is acceptable
- [ ] Multi-line text works
- [ ] Special characters work
- [ ] Long text wraps properly
- [ ] Multiple print jobs work
- [ ] Error handling is robust
- [ ] Disconnection is clean
- [ ] No memory leaks
- [ ] Performance is acceptable
- [ ] Code is documented
- [ ] Tests are written

---

## 🚀 Quick Start for BLE Development

1. **Setup test environment:**

   ```bash
   export PRINTER_TYPE=ble
   export PRINTER_DEVICE_NAME=Peripage
   ```

2. **Enable verbose logging:**
   Add debug logs in `ble.go` to trace execution

3. **Test incrementally:**

   - First: Just discovery
   - Then: Connection
   - Then: Handshake
   - Finally: Printing

4. **Use mock printer for API testing:**
   While developing BLE, keep API tests working with mock

---

## 📝 Notes

- Keep existing mock printer working during BLE development
- Add feature flags if needed to enable/disable incomplete features
- Document all findings in code comments
- Update this TODO as you progress
- Ask for help in thermal printer communities if stuck

---

## 🎯 Success Criteria

The BLE implementation is complete when:

1. ✅ Can automatically discover and connect to Peripage A6
2. ✅ Can print plain text reliably
3. ✅ Can print formatted JSON
4. ✅ Can handle connection errors gracefully
5. ✅ Can reconnect after printer power cycle
6. ✅ Print quality is readable and acceptable
7. ✅ Code is production-ready with error handling
8. ✅ All TODOs in code are resolved

---

**Good luck with the implementation! 🎉**
