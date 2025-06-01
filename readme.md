# Bamos

[![Go Report Card](https://goreportcard.com/badge/github.com/mayloo89/bamos)](https://goreportcard.com/report/github.com/mayloo89/bamos)
[![Build Status](https://github.com/mayloo89/bamos/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/mayloo89/bamos/actions)
[![codecov](https://codecov.io/gh/mayloo89/bamos/branch/main/graph/badge.svg)](https://codecov.io/gh/mayloo89/bamos)

Bamos is a web application built in Go that provides information and tools for public transport in the Metropolitan Area of Buenos Aires, Argentina, using the [CABA Transport API](https://api-transporte.buenosaires.gob.ar/console).

## Features
- Search for bus lines and view route information
- Check allowed parking rules for a given location
- View real-time vehicle positions (GTFS)
- Responsive web UI with Bootstrap
- Session management and CSRF protection

## Tech Stack
- **Language:** Go 1.22.5
- **Web Framework:** [chi router](https://github.com/go-chi/chi/v5)
- **Session Management:** [scs](https://github.com/alexedwards/scs/v2)
- **CSRF Protection:** [nosurf](https://github.com/justinas/nosurf)
- **Testing:** [testify](https://github.com/stretchr/testify)

## Project Structure
```
cmd/bamos/         # Main application entrypoint
internal/          # Application logic (handlers, services, helpers, forms, etc.)
  handler/         # HTTP handlers
  services/        # API clients and business logic
  helpers/         # Error handling and utilities
  forms/           # Form validation
  model/           # Template data models
  render/          # Template rendering
  config/          # App configuration
  ...
static/            # Static assets (images, routes info)
templates/         # HTML templates
migrations/        # Database migrations
```

## Setup & Usage
1. **Clone the repository:**
   ```sh
   git clone https://github.com/mayloo89/bamos.git
   cd bamos
   ```
2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **Set environment variables:**
   - (Recommended) Move sensitive data like API keys to environment variables or a `.env` file.
4. **Run the application:**
   ```sh
   go run ./cmd/bamos
   ```
5. **Access the app:**
   Open [http://localhost:8080](http://localhost:8080) in your browser.

## API Endpoints
- `GET /` — Home page
- `GET /colectivos/search` — Search bus lines
- `POST /colectivos/search` — Search bus lines (form submit)
- `GET /colectivos/vehiclePositionsSimple` — View vehicle positions
- `GET /transit/allowed-parking` — Allowed parking form
- `POST /transit/allowed-parking` — Query allowed parking rules

## Development & Testing
- **Run tests:**
  ```sh
  go test ./...
  ```
- **Test coverage:**
  Use `go test -cover ./...` to check code coverage.
- **Mocking:**
  External API calls are mocked in tests for reliability and speed.

## Best Practices Followed
- Dependency injection for services and HTTP clients
- Separation of concerns (handlers, services, helpers, etc.)
- No external API dependencies in tests (all are mocked)
- Error handling and logging
- Secure session and CSRF management
- Environment variable support for sensitive data

## Contributing
Pull requests and issues are welcome! Please open an issue to discuss your ideas or report bugs.

## Author

Made with ❤️ by [Sebastian S.](https://github.com/mayloo89)

## License

MIT License — see [LICENSE](LICENSE) for details.  
This project is open source under the MIT License. You can use, modify, and distribute it freely.

