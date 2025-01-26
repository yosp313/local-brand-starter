# AI Content Creation Backend

This is the backend service for the AI Content Creation SaaS tool, built with Go, Gin, and Cloudflare AI.

## Features

- AI-powered content generation using Cloudflare AI Workers
- Support for both Mistral-7B and Llama2-7B models
- Real-time content streaming via Server-Sent Events (SSE)
- PostgreSQL database for user and content management
- Token-based rate limiting
- CORS support for frontend integration

## Prerequisites

- Go 1.21 or later
- sqlite
- Cloudflare account with AI Workers access

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a PostgreSQL database:
```bash
createdb ai_content_creation
```

4. Copy the example environment file and update it with your credentials:
```bash
cp .env.example .env
```

Update the following variables in `.env`:
- `DATABASE_URL`: Your PostgreSQL connection string
- `CLOUDFLARE_ACCOUNT_ID`: Your Cloudflare account ID
- `CLOUDFLARE_API_TOKEN`: Your Cloudflare API token
- `FRONTEND_URL`: Your frontend application URL

5. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### Content Generation
```
POST /api/v1/generate
Content-Type: application/json

{
    "model": "mistral-7b",
    "prompt": "Write a blog post about AI"
}
```

### User Management
```
POST /api/v1/users
GET /api/v1/users/:id
```

## Development

To run the server in development mode with hot reload:
```bash
go install github.com/cosmtrek/air@latest
air
```

## Testing

Run the tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
