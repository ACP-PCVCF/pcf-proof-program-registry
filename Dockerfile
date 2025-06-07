FROM golang:1.23

WORKDIR /app

# Copy go.mod and go.sum separately to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary (better than running go run in container)
RUN go build -o app main.go

EXPOSE 8080

# Run the compiled binary instead of `go run`
CMD ["./app"]