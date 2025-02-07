# TalentLensService

A Go-based microservice for handling user authentication, metrics tracking, and backend interactions for the talent management platform.

## Features

- User authentication and authorization
- JWT-based session management
- Metrics collection and analysis
- RESTful API endpoints
- CORS support

## Prerequisites

- Go 1.23 or higher

## Installation

```bash
git clone https://github.com/edvardsanta/TalentLensService
cd TalentLensService
go mod download
```

## Configuration

Create a `.env` file in the project root:

```env
DB_CONNECTION_STRING=
DB_DRIVER=
JWT_SECRET_KEY=
ALLOWED_ORIGINS=*
SERVICE_NAME=
SERVICE_VERSION=
OTEL_SDK_DISABLED=true
```

## Running the Service

Development:
```bash
go run cmd/api/main.go
```

Production:
```bash
go build -o platform-service cmd/api/main.go
./platform-service
```

Docker:
```bash
docker build -t platform-service .
docker run -p 8080:8080 platform-service
```