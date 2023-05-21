VERSION 0.7
FROM golang:1.20
WORKDIR /pom

deps:
    COPY go.mod go.sum ./
    COPY internal ./internal
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build-api:
    FROM +deps
    COPY cmd/api ./cmd/api
    RUN go build -o build/api ./cmd/api/main.go
    SAVE ARTIFACT build/api AS LOCAL build/api

build-vault:
    FROM +deps
    COPY cmd/vault ./cmd/vault
    RUN go build -o build/vault ./cmd/vault/main.go
    SAVE ARTIFACT build/vault AS LOCAL build/vault

build-all:
    BUILD +build-api
    BUILD +build-vault

docker-api:
    COPY +build-api/api .
    ENTRYPOINT ["/pom/api"]
    SAVE IMAGE ghcr.io/madebaruna/pom-moe-api:latest

docker-vault:
    COPY +build-vault/vault .
    ENTRYPOINT ["/pom/vault"]
    SAVE IMAGE ghcr.io/madebaruna/pom-moe-vault:latest

docker-all:
    BUILD +docker-api
    BUILD +docker-vault