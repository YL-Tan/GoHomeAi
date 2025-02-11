# Use official lightweight Go image
FROM golang:1.21 as builder

WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o main

# Run the compiled binary
CMD ["/app/main"]
