# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder


LABEL authors="marston"

COPY . "/go/src/github.com/TheMarstonConnell/musicapi"
COPY "/etc/ssl/certs/ca-certificates.crt" "/etc/ssl/certs/ca-certificates.crt"
#COPY "cert.pem" "/etc/ssl/cert.pem"

WORKDIR "/go/src/github.com/TheMarstonConnell/musicapi"

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /musicapi


EXPOSE 9797
CMD ["/musicapi"]

############################

FROM scratch

COPY --from=builder /musicapi .

EXPOSE 8000

CMD ["/musicapi"]
