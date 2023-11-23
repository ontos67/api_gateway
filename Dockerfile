FROM golang:1.15 AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/api_gateway.go


FROM debian:stretch

COPY --from=builder /code/api_gateway /usr/local/bin/api_gateway

RUN chmod +x /usr/local/bin/api_gateway

ENTRYPOINT [ "api_gateway" ]