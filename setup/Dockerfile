# Use the official Go image to provide a build environment

# Use the official Go image
FROM golang:1.17

WORKDIR /app

# Copy your code into the container
COPY . .

# Download dependencies
RUN go mod download

# Run the publisher application when the Docker image is run
CMD ["go", "run", "main.go"]

