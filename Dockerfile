FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN go build -o main .

# Use an absolute path for the CMD
CMD ["/app/main"]
