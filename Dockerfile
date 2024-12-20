FROM golang:1.22-alpine

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /main cmd/server/main.go

EXPOSE 8888

CMD ["/main"]