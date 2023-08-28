FROM golang:1.20.7-alpine3.18 as builder

WORKDIR build

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/build/main /main

USER nobody

ENTRYPOINT ["/main"]