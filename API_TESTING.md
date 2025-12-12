# API Testing Guide

Quick reference for testing the Peripage Printer API.

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Health Check

Check if the service is running.

**Request:**

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{
  "status": "healthy"
}
```

---

### 2. Print Plain Text

Send plain text to the printer.

**Request:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello from Peripage!"
  }'
```

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
  "error": "text cannot be empty"
}
```

---

### 3. Print JSON Data

Send JSON data to be pretty-printed on the printer.

**Request:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "name": "John Doe",
      "email": "john@example.com",
      "order_id": 12345,
      "items": [
        {"product": "Widget A", "qty": 2},
        {"product": "Widget B", "qty": 1}
      ]
    }
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Print job completed successfully"
}
```

**What gets printed:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "order_id": 12345,
  "items": [
    {
      "product": "Widget A",
      "qty": 2
    },
    {
      "product": "Widget B",
      "qty": 1
    }
  ]
}
```

---

### 4. Print Text and JSON (Text is ignored if data exists)

If both `text` and `data` are provided, only `data` is printed.

**Request:**

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "text": "This will be ignored",
    "data": {
      "message": "Only this JSON will be printed"
    }
  }'
```

---

## Test Scenarios

### Scenario 1: Receipt Printing

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "receipt_id": "RCP-2025-001",
      "date": "2025-12-12",
      "customer": "Jane Smith",
      "items": [
        {
          "item": "Coffee",
          "price": 3.50,
          "qty": 2
        },
        {
          "item": "Croissant",
          "price": 2.75,
          "qty": 1
        }
      ],
      "subtotal": 9.75,
      "tax": 0.98,
      "total": 10.73
    }
  }'
```

### Scenario 2: Label Printing

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "text": "FRAGILE\nHandle with Care\nShip to: New York\nOrder #12345"
  }'
```

### Scenario 3: Status Update

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "status": "Order Confirmed",
      "order_number": "ORD-2025-12-12-001",
      "estimated_delivery": "Dec 15, 2025",
      "tracking": "1Z999AA10123456784"
    }
  }'
```

---

## Error Testing

### Empty Request

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**

```json
{
  "error": "Either 'text' or 'data' must be provided"
}
```

### Invalid JSON

```bash
curl -X POST http://localhost:8080/print \
  -H "Content-Type: application/json" \
  -d 'not valid json'
```

**Response:**

```json
{
  "error": "Invalid request body: ..."
}
```

---

## Using HTTPie (Alternative)

If you have HTTPie installed:

```bash
# Health check
http GET http://localhost:8080/health

# Print text
http POST http://localhost:8080/print text="Hello World"

# Print JSON
http POST http://localhost:8080/print data:='{"key":"value"}'
```

---

## Swagger UI

For interactive API testing, visit:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:

- Interactive endpoint testing
- Request/response examples
- Schema documentation
- Try-it-out functionality

---

## Testing with Postman

1. Import the Swagger JSON:

   ```
   http://localhost:8080/swagger/doc.json
   ```

2. Or manually create requests:
   - Method: POST
   - URL: `http://localhost:8080/print`
   - Headers: `Content-Type: application/json`
   - Body: JSON payload

---

## Mock Printer Output

When using `PRINTER_TYPE=mock`, you'll see output in the server logs:

```
[peripage] 2025/12/12 10:00:00 === MOCK PRINTER OUTPUT ===
Hello from Peripage!
[peripage] 2025/12/12 10:00:00 === END MOCK PRINTER OUTPUT ===
```

This is perfect for testing without hardware!
