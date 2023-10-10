FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app
COPY src .

RUN go install

FROM gcr.io/distroless/static-debian11:latest

LABEL maintainer="opensource@breadnet.co.uk"
LABEL org.opencontainers.image.source="https://github.com/userbradley/prom-status"
LABEL org.opencontainers.image.description="Polls HTTP endpoints and reports latency via Prometheus"
LABEL org.opencontainers.image.base.name="gcr.io/distroless/static-debian11:latest"

COPY --from=builder /go/bin/prom-status /usr/local/bin/prom-status

WORKDIR /app

CMD ["prom-status"]
