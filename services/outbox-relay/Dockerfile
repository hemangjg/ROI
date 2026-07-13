FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.work go.work.sum ./
COPY packages ./packages
COPY db ./db
COPY services/outbox-relay ./services/outbox-relay

WORKDIR /src/services/outbox-relay
ENV GOWORK=off
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/relay ./cmd/relay

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /bin/relay /usr/local/bin/relay
COPY services/outbox-relay/config/defaults.yaml /app/config/defaults.yaml
EXPOSE 8086
ENTRYPOINT ["/usr/local/bin/relay"]