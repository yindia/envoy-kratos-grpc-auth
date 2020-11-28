FROM golang:alpine AS builder

RUN apk --no-cache add make
COPY . /app

RUN make -C /app

FROM alpine

COPY --from=builder /app/server /app/server

CMD ["/app/server", "-users", "/etc/users.json"]
