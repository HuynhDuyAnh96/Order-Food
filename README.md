# Dish API Service

A RESTful API service for managing Vietnamese dishes with filtering, pagination, and featured dish functionality.

## Features

- **GET /api/dishes** - Retrieve dishes with advanced filtering and pagination
- **GET /api/dishes/featured** - Get 2 featured popular dishes for homepage
- **GET /health** - Health check endpoint

## API Endpoints

### 1. Get All Dishes
```
GET /api/dishes
```

#### Query Parameters:
- `page` (number): Current page number (default: 1)
- `limit` (number): Items per page (default: 10, max: 100)
- `category` (string): Filter by category (e.g., "noodles", "sandwich", "rice")
- `cooking_method` (string): Filter by cooking method (e.g., "grilled", "boiled", "fried")
- `is_popular` (boolean): Filter popular dishes only
- `sort` (string): Sort order
  - `price_asc`: Price ascending
  - `price_desc`: Price descending  
  - `rating_desc`: Rating descending
  - Default: Sort by name

#### Response Format:
```json
{
  "success": true,
  "data": {
    "dishes": [
      {
        "id": "1",
        "name": "Phở Bò",
        "description": "Traditional Vietnamese beef noodle soup",
        "price": 85000,
        "category": "noodles",
        "cooking_method": "boiled",
        "is_popular": true,
        "rating": 4.8,
        "image_url": "https://example.com/pho-bo.jpg",
        "created_at": "2025-08-17T14:52:35.251279+07:00",
        "updated_at": "2025-09-17T14:52:35.251385+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 8,
      "totalPages": 1
    }
  }
}
```

### 2. Get Featured Dishes
```
GET /api/dishes/featured
```

Returns 2 most popular dishes sorted by rating for homepage display.

#### Response Format:
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "name": "Phở Bò",
      "description": "Traditional Vietnamese beef noodle soup",
      "price": 85000,
      "category": "noodles",
      "cooking_method": "boiled",
      "is_popular": true,
      "rating": 4.8,
      "image_url": "https://example.com/pho-bo.jpg",
      "created_at": "2025-08-17T14:52:35.251279+07:00",
      "updated_at": "2025-09-17T14:52:35.251385+07:00"
    }
  ]
}
```

## Example Usage

### Get all dishes with pagination:
```bash
curl "http://localhost:8080/api/dishes?page=1&limit=5"
```

### Filter by category and sort by price:
```bash
curl "http://localhost:8080/api/dishes?category=noodles&sort=price_asc"
```

### Get only popular dishes:
```bash
curl "http://localhost:8080/api/dishes?is_popular=true&sort=rating_desc"
```

### Filter by cooking method:
```bash
curl "http://localhost:8080/api/dishes?cooking_method=grilled"
```

### Get featured dishes for homepage:
```bash
curl "http://localhost:8080/api/dishes/featured"
```

## Running the Service

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## Sample Data

The service includes 8 sample Vietnamese dishes:
- Phở Bò (Popular)
- Bánh Mì Thịt Nướng (Popular)
- Bún Chả
- Cơm Tấm
- Gỏi Cuốn (Popular)
- Chả Cá Lã Vọng
- Bánh Xèo
- Bún Bò Huế (Popular)

## CORS Support

The API includes CORS headers for frontend integration:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Architecture

The project follows clean architecture principles:
- **Domain Layer**: Business entities and models
- **Application Layer**: Business logic and services
- **Infrastructure Layer**: HTTP handlers, repositories, and external dependencies
- **Main**: Application entry point and dependency injection
