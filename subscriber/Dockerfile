FROM golang:1.17
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Expose port 8080 to the outside
EXPOSE 8080
CMD ["/app/main"]

