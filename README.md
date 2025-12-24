# Airline Search Simulator

A high-performance flight search aggregator that queries multiple airline providers in parallel and returns unified search results with filtering, sorting, and ranking capabilities.

## Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

## Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd airline-search-simulator
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
go mod vendor

#Or use
make install-deps
```

### 3. Configure the Application

Copy the example configuration file:

```bash
cp .env.yaml.example .env.yaml
```

Edit `.env.yaml` to customize settings (optional):

```yaml
server:
  port: 8080
  timeout: "30s"

provider:
  timeout: "5s"
  providers:
    garuda:
      enabled: true
      name: "Garuda Indonesia"
      data_path: "test_data/garuda_indonesia_search_response.json"
      response_time: "500ms"
      failure_rate: 0.0
    # ... other providers
```

### 4. Build the Application

```bash
# Using Make
make build

# Or directly with Go
go build -o bin/server ./cmd/server
```

## Running the Application

### Start the Server

```bash
# Using Make
make run

# Or directly
./bin/server

# Or with go run
go run ./cmd/server/main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env.yaml`).

## API Endpoints

### 1. Health Check

Check if the service is running:

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. List Providers

Get list of available airline providers:

```bash
curl http://localhost:8080/providers
```

**Response:**
```json
{
  "providers": [
    "Garuda Indonesia",
    "Lion Air",
    "Batik Air",
    "AirAsia"
  ]
}
```

### 3. Search Flights

Search for flights with various filters and options.

#### Basic Search

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy"
  }'
```

#### Search with Filters

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 2,
    "cabinClass": "economy",
    "filters": {
      "maxPrice": 1500000,
      "maxStops": 0,
      "airlines": ["Garuda Indonesia", "Lion Air"],
      "departureTime": {
        "start": 6,
        "end": 12
      },
      "maxDuration": 180
    }
  }'
```

#### Search with Sorting

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy",
    "sortBy": "price",
    "sortOrder": "asc"
  }'
```

#### Complete Search Example

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "returnDate": "2025-12-22",
    "passengers": 2,
    "cabinClass": "economy",
    "filters": {
      "minPrice": 500000,
      "maxPrice": 2000000,
      "maxStops": 1,
      "airlines": ["Garuda Indonesia", "Lion Air"],
      "departureTime": {
        "start": 6,
        "end": 18
      },
      "arrivalTime": {
        "start": 8,
        "end": 22
      },
      "maxDuration": 180
    },
    "sortBy": "price",
    "sortOrder": "asc",
    "returnFilters": {
        "maxPrice": 650000
      },
      "returnSortBy": "",
      "returnSortOrder": ""
  }'
```

**Response Format:**
```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "cabin_class": "economy"
  },
  "metadata": {
    "total_results": 15,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 285,
    "cache_hit": false
  },
  "flights": [
    {
      "id": "QZ7250_AirAsia",
      "provider": "AirAsia",
      "airline": {
        "name": "AirAsia",
        "code": "QZ"
      },
      "flight_number": "QZ7250",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T15:15:00+07:00",
        "timestamp": 1734246900
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T20:35:00+08:00",
        "timestamp": 1734267300
      },
      "duration": {
        "total_minutes": 260,
        "formatted": "4h 20m"
      },
      "stops": 1,
      "price": {
        "amount": 485000,
        "currency": "IDR"
      },
      "available_seats": 88,
      "cabin_class": "economy",
      "amenities": [],
      "baggage": {
        "carry_on": "Cabin baggage only",
        "checked": "Additional fee"
      }
    }
  ]
}
```

## Request Parameters

### Required Fields

- `origin` (string): Departure airport code (3-letter IATA code)
- `destination` (string): Arrival airport code (3-letter IATA code)
- `departureDate` (string): Departure date (YYYY-MM-DD format)
- `passengers` (int): Number of passengers (1-9)
- `cabinClass` (string): Cabin class (`economy`, `premium`, `business`, `first`)

### Optional Fields

- `returnDate` (string): Return date for round-trip flights
- `sortBy` (string): Sort field (`price`, `duration`, `departure`, `arrival`, `stops`)
- `sortOrder` (string): Sort order (`asc` or `desc`)

### Filter Options

- `minPrice` (float): Minimum price in IDR
- `maxPrice` (float): Maximum price in IDR
- `maxStops` (int): Maximum number of stops
- `airlines` (array): Array of airline names (case-insensitive)
- `departureTime.start` (int): Earliest departure hour (0-23)
- `departureTime.end` (int): Latest departure hour (0-23)
- `arrivalTime.start` (int): Earliest arrival hour (0-23)
- `arrivalTime.end` (int): Latest arrival hour (0-23)
- `maxDuration` (int): Maximum flight duration in minutes

## Common Airport Codes

| Code | City |
|------|------|
| CGK | Jakarta |
| DPS | Denpasar (Bali) |
| SUB | Surabaya |
| JOG | Yogyakarta |
| UPG | Makassar |
| KNO | Medan |

## Development


## Features

- **Parallel Provider Queries**: Queries multiple airline providers simultaneously
- **Intelligent Caching**: Caches search results to improve performance
- **Advanced Filtering**: Filter by price, stops, airlines, departure/arrival times, and duration
- **Flexible Sorting**: Sort results by price, duration, departure time, or number of stops
- **Smart Ranking**: Automatically scores and ranks flights based on multiple factors
- **Provider Filtering**: Only queries relevant providers when airline filter is specified
- **Error Handling**: Graceful error handling with partial results support
- **Validation**: Comprehensive input validation for all request parameters