FROM golang:1.17
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "main.go"]

