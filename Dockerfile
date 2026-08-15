FROM golang:1.25-alpine

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum .

# Install dependencies
RUN go mod tidy

# Copy source code
COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]
