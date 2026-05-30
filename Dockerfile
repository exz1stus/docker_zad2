# ETAP 1: Budowanie
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

LABEL org.opencontainers.image.authors="Denys Petrov"

# Instalacja narzędzi pomocniczych
RUN apk add --no-cache upx ca-certificates

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go ./

# Kompilacja z uwzględnieniem docelowej architektury (TARGETOS i TARGETARCH dostarczane przez buildx)
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -a -installsuffix cgo -o weather-app .

RUN upx --best weather-app && chmod +x weather-app

# ETAP 2: Obraz końcowy
FROM scratch

LABEL org.opencontainers.image.authors="Denys Petrov"

COPY --from=builder /app/weather-app /weather-app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY data.json /data.json

HEALTHCHECK --interval=30s --timeout=3s CMD ["/weather-app", "health"]

EXPOSE 8080
ENTRYPOINT ["/weather-app"]