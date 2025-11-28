FROM golang:1.25 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .
COPY frontend/. ./frontend/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /dnsmasq-k8s ./cmd/api

FROM debian:stable-slim

RUN apt-get update && apt-get install -y dnsmasq supervisor procps && apt-get clean

COPY --from=builder /dnsmasq-k8s /dnsmasq-k8s

WORKDIR /

COPY --from=builder /app/frontend/src /frontend/src

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

ENV DNS_ENABLED=false
ENV DHCP_ENABLED=false

EXPOSE 8080
EXPOSE 53 53/udp
EXPOSE 67 67/udp

CMD ["/usr/bin/supervisord"]