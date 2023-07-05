# syntax=docker/dockerfile:1

FROM golang:1.20


# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-server

EXPOSE 5000
# Run

CMD ["./go-server"]