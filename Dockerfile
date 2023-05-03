FROM golang:1.20.1-alpine as builder

WORKDIR /usr/src/goctopus

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/goctopus ./cmd/goctopus/goctopus.go

FROM alpine:3.14
COPY --from=builder /usr/local/bin/goctopus ./goctopus
ENTRYPOINT [ "./goctopus" ]