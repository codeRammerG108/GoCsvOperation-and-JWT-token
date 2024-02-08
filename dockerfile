# Use a minimal base image with Go runtime
FROM golang:1.21.5-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Copy the Go application source code to the container
COPY . .

# Build the Go application (e.g., go build -o myapp)
RUN go build -o myapp

# Expose the port your Go application will run on
EXPOSE 8080

# Define the command to run your Go application
CMD ["./myapp"]
