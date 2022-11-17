FROM golang:1.18-alpine3.14 AS builder

WORKDIR /build
COPY . /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o shortener_svc cmd/main.go

FROM scratch

WORKDIR /shortener
COPY --from=builder ["/build/shortener_svc", "/shortener/"]

EXPOSE 7002
ENTRYPOINT ["/shortener/shortener_svc"]
