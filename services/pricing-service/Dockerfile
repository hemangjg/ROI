FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.work go.work.sum ./
COPY packages ./packages
COPY db ./db
COPY services/pricing-service ./services/pricing-service

WORKDIR /src/services/pricing-service
ENV GOWORK=off
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /bin/server /usr/local/bin/server
COPY services/pricing-service/config/defaults.yaml /app/config/defaults.yaml
EXPOSE 8084 9092
ENTRYPOINT ["/usr/local/bin/server"]