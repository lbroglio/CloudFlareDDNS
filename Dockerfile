FROM docker.io/library/golang:latest as builder

WORKDIR /app
COPY . .

RUN make test build-static

FROM scratch

COPY --from=builder /app/bin/cloudflareddns /usr/local/bin/cloudflareddns

CMD ["cloudflareddns"]