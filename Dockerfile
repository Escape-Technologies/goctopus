FROM golang:1.20.1-alpine

WORKDIR /usr/src/goctopus

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/goctopus ./cmd/goctopus/goctopus.go

ENTRYPOINT [ "goctopus" ]