FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o opiso ./cmd/opiso/opiso.go
CMD ["./opiso"]
