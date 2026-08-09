FROM docker.io/library/golang:latest as builder

COPY . .

RUN make test build

FROM docker.io/library/alpine:latest

COPY --from=builder /go/bin/cloudflareddns /usr/local/bin/cloudflareddns

CMD ["cloudflareddns"]