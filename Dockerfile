# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN go mod download

# Install godotenv to load environment variables from .env
RUN go get github.com/joho/godotenv

# Copy .env file into the Docker image
COPY .env .env

# Build the Go application
RUN go build -o main .

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the executable with godotenv and pass the connection string as an argument
CMD ["./main", "-mongouri", "${MONGOURI}"]
