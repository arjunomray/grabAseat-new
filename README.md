# Grab A Seat

A modern event ticketing platform built with Go and React, allowing users to create, manage, and book event tickets with ease.

## Tech Stack

### Backend
- Go 1.23.4
- Gin Web Framework
- GORM with PostgreSQL
- JWT Authentication

### Frontend
- React 18
- Vite
- Tailwind CSS
- React Router DOM
- QRCode Generation

### Infrastructure
- Docker & Docker Compose
- PostgreSQL 16

## Getting Started

### Prerequisites
- Docker and Docker Compose installed on your machine

### Quick Start

1. Clone the repository:

```bash
git clone <repository-url>
cd grab-a-seat
```

2. Start the application:
```bash
docker-compose up --build
```

3. Access the application:
- Frontend: http://localhost:5173
- API: http://localhost:8080

### Available Endpoints

Public Routes:
- `POST /register` - Create new account
- `POST /login` - User authentication

Protected Routes (requires JWT):
- `POST /events` - Create event
- `GET /events` - List all events
- `GET /events/:id` - Get event details
- `POST /tickets` - Book ticket
- `GET /tickets` - View all tickets
- `GET /tickets/:id` - View user tickets

## Development

### Project Structure
```
.
├── backend/          # Go backend service
├── frontend/         # React frontend
├── docker-compose.yml
├── Dockerfile.backend
└── Dockerfile.frontend
```

### Environment Setup
Create a `.env` file:
```env
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=eventer
DB_PORT=5432
```

### Stopping the Application
```bash
docker-compose down
```
