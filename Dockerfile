# Use the official Golang image as the base image
FROM golang:1.24.1 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o image-renamer-admission-plugin ./cmd/image-renamer-admission-plugin

# Use a minimal base image for the final stage
FROM gcr.io/distroless/base-debian12

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/image-renamer-admission-plugin .

# Copy the configuration file
COPY config.yaml .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./image-renamer-admission-plugin"]
