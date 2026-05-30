github repo: https://github.com/exz1stus/docker_zad1.git

a. zbudowania opracowanego obrazu kontenera,

    docker build -t weather-app:v1 .

b. uruchomienia kontenera na podstawie zbudowanego obrazu,

    docker run -d -p 8080:8080 --name weatherapp weather-app:v1

c. sposobu uzyskania informacji z logów, które wygenerowałą opracowana aplikacja podczas
uruchamiana kontenera (patrz: punkt 1a),

    docker logs weatherapp

        Data uruchomienia: 2026-04-30 17:55:27
        Autor: Denys Petrov
        Port TCP: 8080

d. sprawdzenia, ile warstw posiada zbudowany obraz oraz jaki jest rozmiar obrazu.

    1. docker images weather-app:v1

    IMAGE            ID             DISK USAGE   CONTENT SIZE   EXTRA
    weather-app:v1   2ad388596e76        4.9MB         2.36MB

    2. docker history weather-app:v1

    IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
    2ad388596e76   14 seconds ago   ENTRYPOINT ["/weather-app"]                     0B        buildkit.dockerfile.v0
    <missing>      14 seconds ago   EXPOSE [8080/tcp]                               0B        buildkit.dockerfile.v0
    <missing>      14 seconds ago   HEALTHCHECK &{["CMD" "/weather-app" "health"…   0B        buildkit.dockerfile.v0
    <missing>      14 seconds ago   LABEL org.opencontainers.image.authors=Denys…   0B        buildkit.dockerfile.v0
    <missing>      14 seconds ago   COPY /etc/ssl/certs/ca-certificates.crt /etc…   242kB     buildkit.dockerfile.v0
    <missing>      14 seconds ago   COPY data.json /data.json # buildkit            8.19kB    buildkit.dockerfile.v0
    <missing>      14 seconds ago   COPY /app/weather-app /weather-app # buildkit   2.29MB    buildkit.dockerfile.v0

Czesc nieobowiazkowa

tworzenie buildera

    docker buildx create --name zad1_builder --driver docker-container --use

    docker buildx inspect --bootstrap

docker scout cves mbmbpididi/weather-app:latest

❯ docker scout cves mbmbpididi/weather-app:latest
✓ Pulled
✓ Image stored for indexing
✓ Indexed 0 packages
✓ Provenance obtained from attestation
✓ No vulnerable package detected

## Overview

                   │               Analyzed Image

───────────────────┼─────────────────────────────────────────────
Target │ mbmbpididi/weather-app:latest  
 digest │ a5070d23ef8c  
 platform │ linux/amd64  
 provenance │ https://github.com/exz1stus/docker_zad1.git
│ bcc31abcaf6c706476469b7c86dace90df53fac6  
 vulnerabilities │ 0C 0H 0M 0L  
 size │ 2.4 MB  
 packages │ 0

## Packages and Vulnerabilities

No vulnerable packages detected

Użycie obrazu bazowego scratch zazwyczaj skutkuje zerową liczbą podatności

Obraz został poddany analizie pod kątem podatności (CVE). Wynik analizy wskazuje na brak jakichkolwiek zagrożeń (0 Critical, 0 High, 0 Medium, 0 Low).

        ❯ docker buildx imagetools inspect mbmbpididi/weather-app:latest
        Name: docker.io/mbmbpididi/weather-app:latest
        MediaType: application/vnd.oci.image.index.v1+json
        Digest: sha256:0716cb6f9fd1904615794d0eaa2e97940b98bd5fdf854e553f6f5fe4cc56fe88

        Manifests:
        Name: docker.io/mbmbpididi/weather-app:latest@sha256:a5070d23ef8cd6f4e29fd3f3fa398804c2dba30ffc96260d6430d3f1a9c693bc
        MediaType: application/vnd.oci.image.manifest.v1+json
        Platform: linux/amd64

        Name: docker.io/mbmbpididi/weather-app:latest@sha256:204e9f9f03cb9714711c5b565d45cae1e95ef9e00871391b4bf5725874c7218f
        MediaType: application/vnd.oci.image.manifest.v1+json
        Platform: linux/arm64

        Name: docker.io/mbmbpididi/weather-app:latest@sha256:cef052efbf8a382e9a2760881761c8c2a731b4ad9a97ab24862d966a68064c66
        MediaType: application/vnd.oci.image.manifest.v1+json
        Platform: unknown/unknown
        Annotations:
        vnd.docker.reference.digest: sha256:a5070d23ef8cd6f4e29fd3f3fa398804c2dba30ffc96260d6430d3f1a9c693bc
        vnd.docker.reference.type: attestation-manifest

        Name: docker.io/mbmbpididi/weather-app:latest@sha256:414155c840281aad652483f3e4e13e767ee3b0e063f97740edfe2036f01182b6
        MediaType: application/vnd.oci.image.manifest.v1+json
        Platform: unknown/unknown
        Annotations:
        vnd.docker.reference.digest: sha256:204e9f9f03cb9714711c5b565d45cae1e95ef9e00871391b4bf5725874c7218f
        vnd.docker.reference.type: attestation-manifest

potwierdza poprawną strukturę Multi-platform:

linux/amd64
linux/arm64
