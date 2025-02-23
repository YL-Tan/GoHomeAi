FROM golang:1.23.2 AS base
 
WORKDIR /usr/src/app
 
COPY go.mod go.sum ./

RUN go mod download

COPY . .
 
ENV APP_ENV=development
 
# ===========================
# Development Stage
# ===========================
FROM base AS development
 
WORKDIR /usr/src/app

COPY --from=base /usr/src/app /usr/src/app

RUN go install github.com/cosmtrek/air@latest  # Hot-reload tool

CMD ["air"]  # Auto-reload on code changes
 
# ===========================
# Build Stage
# ===========================
FROM base AS build
 
# Build the application for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o /app main.go

# ===========================
# Staging Stage
# ===========================
FROM golang:1.23.2 AS staging
 
# Set the working directory
WORKDIR /usr/src/app
 
COPY go.mod go.sum ./

RUN go mod download

COPY . .
 
# Use .env file for staging if necessary
ENV APP_ENV=staging
 
# Expose the staging port
EXPOSE ${PORT}
 
# Run the application in staging mode
CMD ["./app"]

# ===========================
# Production Stage
# ===========================
FROM golang:1.23.2 AS production
 
# Set the working directory
WORKDIR /usr/src/app
 
COPY go.mod go.sum ./

RUN go mod download

COPY . .
 
# Use .env file for production if necessary
ENV APP_ENV=production
 
# Expose the production port
EXPOSE ${PORT}
 
# Run the application in production mode
CMD ["./app"]