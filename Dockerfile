FROM golang:1.24-alpine AS BUILDER
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o server main.go

FROM alpine
WORKDIR /app

COPY --from=BUILDER /app/server .
COPY --from=BUILDER /app/migrations ./migrations

EXPOSE 8080
CMD ["./server"]