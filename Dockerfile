FROM golang:1.14.6-alpine3.12 as build

RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN mkdir /app
WORKDIR /app

# Trick for caching the dependencies better based on changes in the go.mod and go.sum files
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app
ENV CGO_ENABLED=0
RUN go build -o bin/k8s-metadata-injection cmd/server/main.go

FROM alpine:3.12.0

RUN mkdir /app
RUN apk add --update openssl
COPY entrypoint.sh /app
COPY --from=build /app/bin/k8s-metadata-injection /app

CMD ["/app/entrypoint.sh"]
