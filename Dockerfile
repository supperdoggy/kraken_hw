# Use the official Golang image to create a build artifact.
FROM golang:1.20-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.mod and go.sum files to the workspace.
COPY go.mod ./

# Download and cache dependencies.
RUN go mod tidy

# Copy the source code to the workspace.
COPY . ./

# Build the Go app.
RUN go build -o main ./main.go

# Use a minimal image.
FROM alpine:3.18

# Copy the binary from the builder image.
COPY --from=builder /app/main /main

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["/main"]
