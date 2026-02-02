
FROM golang:1.24-alpine AS builder

WORKDIR /app


RUN apk add --no-cache make npm

COPY . .

RUN make build-frontend && make build-go-server



FROM alpine:3.22

WORKDIR /app

# result build name is pppk-json in Makefile
COPY --from=builder /app/backend/bin/pppk-json .


ENV PORT=8080

EXPOSE ${PORT}

CMD ["./pppk-json"]
