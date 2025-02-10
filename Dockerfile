FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire source code into the container
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o mongodb-demo .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./mongodb-demo"]