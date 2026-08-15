# User Management API

A REST API for user management and session-based authentication, built with Go, Gin, and MongoDB.

## Features

- User creation, retrieval, pagination, updates, and deletion
- Short-lived JWT access tokens
- Rotating opaque refresh tokens stored in HttpOnly cookies
- MongoDB-backed users and user sessions
- Graceful HTTP server shutdown
- Docker Compose setup for the API, MongoDB, and Mongo Express
- Importable Postman collection with sample requests and responses

## Prerequisites

For the recommended Docker workflow:

- Docker Desktop or Docker Engine with Docker Compose v2

For local development:

- Go 1.25 or newer
- Docker, unless MongoDB is already available separately

## Environment setup

Copy the example configuration before starting the project.

PowerShell:

```powershell
Copy-Item .env.example .env
```

macOS or Linux:

```sh
cp .env.example .env
```

The Compose file also references four variables that are not currently included in `.env.example`. Add them to `.env` with local-only values:

```dotenv
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password
MONGO_EXPRESS_USERNAME=admin
MONGO_EXPRESS_PASSWORD=password
```

Review these important settings:

- `JWT_SECRET`: replace the example value with a long random secret. Never commit a real secret.
- `COOKIE_SECURE=false`: required when testing refresh cookies over local HTTP. Use `true` behind HTTPS.
- `MONGO_URI`: its hostname depends on how the API is run:
  - API running locally: `mongodb://admin:password@localhost:27017`
  - API running in Compose: `mongodb://admin:password@mongodb:27017`

If the MongoDB password contains URL-sensitive characters, percent-encode it in `MONGO_URI`.

## Run with Docker Compose (recommended)

Before starting, verify that `mongodb/initdb.d/mongo-init.js` is a regular JavaScript file. In some checkouts it may be a directory, which causes a Docker bind-mount error. See [Troubleshooting](#troubleshooting) if that is the case.

Set the Compose-compatible MongoDB hostname in `.env`:

```dotenv
MONGO_URI=mongodb://admin:password@mongodb:27017
```

Build and start all services:

```sh
docker compose up --build -d
```

Check their status and follow the API logs:

```sh
docker compose ps
docker compose logs -f app
```

Verify the API after it starts:

PowerShell:

```powershell
Invoke-RestMethod http://localhost:8080/health
```

macOS or Linux:

```sh
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "healthy"
}
```

## Run the API locally

This workflow runs MongoDB and Mongo Express in containers while the Go API runs on the host.

1. Set the localhost MongoDB URI in `.env`:

   ```dotenv
   MONGO_URI=mongodb://admin:password@localhost:27017
   ```

2. Start MongoDB and Mongo Express:

   ```sh
   docker compose up -d mongo mongo-express
   ```

3. Download Go dependencies:

   ```sh
   go mod download
   ```

4. Run the API from the repository root:

   ```sh
   go run ./cmd/main.go
   ```

5. Open `http://localhost:8080/health` to confirm the server is healthy.

Press `Ctrl+C` to stop the local API gracefully.

## Service URLs

| Service | URL |
| --- | --- |
| API | `http://localhost:8080` |
| Health check | `http://localhost:8080/health` |
| Mongo Express | `http://localhost:8081` |

Mongo Express uses `MONGO_EXPRESS_USERNAME` and `MONGO_EXPRESS_PASSWORD` from `.env`.

## API endpoints

All routes use the `/api/v1` prefix. User routes require `Authorization: Bearer <access-token>`.

| Handler | Method | Endpoint | Description |
| --- | --- | --- | --- |
| Auth | `POST` | `/api/v1/auth/login` | Authenticate and set the refresh-token cookie |
| Auth | `POST` | `/api/v1/auth/refresh` | Rotate the refresh token and return a new access token |
| Auth | `POST` | `/api/v1/auth/logout` | Revoke the current session and clear its cookie |
| User | `POST` | `/api/v1/users` | Create a user |
| User | `GET` | `/api/v1/users/:id` | Retrieve a user by UUID |
| User | `GET` | `/api/v1/users?page=1&pageSize=10` | List users with pagination |
| User | `PATCH` | `/api/v1/users/:id` | Update a user's name or email |
| User | `DELETE` | `/api/v1/users/:id` | Delete a user |

The service currently has no public registration endpoint: creating users is protected by bearer authentication. A valid user must therefore already exist in MongoDB before the first login. Provision the initial user through your deployment or database-seeding process, and store its password using the same bcrypt format expected by the service.

For authentication details, see [docs/authentication.md](docs/authentication.md).

## Use the Postman collection

Import [postman/user-mgmt.postman_collection.json](postman/user-mgmt.postman_collection.json) into Postman.

The collection defines:

- `baseUrl`, defaulting to `http://localhost:8080`
- `accessToken`, captured automatically after a successful login or refresh
- `userId`, captured automatically after user creation

Suggested request order:

1. `Auth Handler / Login`
2. `User Handler / Create User`
3. The remaining user requests
4. `Auth Handler / Refresh Token` when needed
5. `Auth Handler / Logout`

Postman's cookie jar manages the HttpOnly refresh-token cookie. The collection contains only synthetic sample credentials and tokens; change the login body to an account provisioned in your database.

## Run tests

Run the complete Go test suite:

```sh
go test ./...
```

At the time this README was written, the repository had pre-existing failures in `service/user_service_test.go` related to password hashing and mock expectations. Other packages passed. Treat failures from the command above as code/test issues to investigate, not necessarily as setup failures.

If Go cannot initialize its build cache on Windows, use a writable task-specific directory:

```powershell
$env:GOCACHE = Join-Path $env:TEMP 'user-mgmt-go-build-cache'
go test ./...
```

## Stop the services

Stop Compose services while keeping MongoDB data:

```sh
docker compose down
```

To also remove the named MongoDB volumes and all stored local data:

```sh
docker compose down --volumes
```

The second command permanently removes the Compose-managed database data.

## Troubleshooting

### MongoDB init-script mount fails

`docker-compose.yml` expects `mongodb/initdb.d/mongo-init.js` to be a file. Check it with:

```powershell
(Get-Item mongodb\initdb.d\mongo-init.js).PSIsContainer
```

The expected result is `False`. If it returns `True`, the path is a directory and must be replaced with the intended MongoDB initialization script from the project source before Compose can mount it correctly.

### API cannot connect to MongoDB

- Use hostname `mongodb` when the API runs inside Compose.
- Use hostname `localhost` when the API runs directly on the host.
- Keep the username and password consistent across `MONGO_URI`, `MONGO_INITDB_ROOT_USERNAME`, and `MONGO_INITDB_ROOT_PASSWORD`.
- Check MongoDB health with `docker compose ps` and logs with `docker compose logs mongo`.

### Refresh or logout returns `401`

- Log in first so Postman or the client receives the refresh cookie.
- Preserve cookies between requests.
- Use `COOKIE_SECURE=false` only for local HTTP; secure cookies are not sent over plain HTTP.
- Keep the refresh cookie path at `/api/v1/auth`.

### Port already in use

The default ports are `8080` for the API, `8081` for Mongo Express, and `27017` for MongoDB. Stop the conflicting service or update both the Compose port mapping and the corresponding local URL/configuration.

### Container does not include recent changes

Rebuild the API image:

```sh
docker compose up --build -d app
```
