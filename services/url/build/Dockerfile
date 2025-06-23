# Stage 1: Build binary
FROM golang:1.24-alpine AS builder
WORKDIR /app

ARG service=user
RUN apk add --no-cache git ca-certificates bash curl tar
COPY . .

# Add private repo auth

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=readonly -o /out/main ./cmd/app/main.go

# Stage 2: Runtime
FROM alpine:3.15
RUN apk --no-cache add ca-certificates bash

ARG env=dev
ARG service
COPY --from=builder /out/main /app/server
COPY configs/config.yaml /configs/
EXPOSE 8081

CMD ["/app/server", "-c", "/configs/config.yaml"]
