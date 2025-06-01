FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make gcc g++ libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server main.go

FROM alpine:latest
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/server .
RUN chmod +x ./server
RUN chown -R appuser ./
USER appuser
RUN ls -la /app

EXPOSE 3000 

CMD ["./server"]
