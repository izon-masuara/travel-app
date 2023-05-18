# Start from the official Golang base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port that the API will listen on
EXPOSE 3000

# Define the command to run the API
CMD ["./main"]
