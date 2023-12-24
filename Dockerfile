# Use an official Go runtime as a base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application files to the container
COPY . .

# Install templ inside the container
RUN go install github.com/a-h/templ/cmd/templ@latest

# Run templ to generate Go files
# RUN templ generate -path="./view"

# Build the Go application
RUN go build -o main cmd/shah/main.go

# Expose the port your application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]