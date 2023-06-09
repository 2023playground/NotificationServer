# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the server files into the container at /app
COPY . /app

# Build the server executablee
RUN go build -o server main.go

# Run the server when the container starts
CMD ["./server"]
