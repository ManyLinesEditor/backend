FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/api cmd/main.go


FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
ENTRYPOINT ["./api"]
