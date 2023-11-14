# NOTE: before building this image, need to modify the template path in main.go:
# from r.LoadHTMLGlob("templates/*") to r.LoadHTMLGlob("/app/templates/*")


# Use the official Golang image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the application source code and the templates folder to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a lightweight Alpine image as the final base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/app .
COPY --from=builder /app/templates /app/templates

# Expose the port the app runs on
EXPOSE 8099

# Command to run the executable
CMD ["./app"]
