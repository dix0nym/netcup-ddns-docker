FROM golang:1.21 as builder

ENV GO111MODULE=on

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o netcup-ddns

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/netcup-ddns /app/netcup-ddns

CMD [ "/app/netcup-ddns" ]