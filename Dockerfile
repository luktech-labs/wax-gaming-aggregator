FROM golang:1.19 as builder
ENV CGO_ENABLED=0

WORKDIR /src
COPY . /src

RUN go build -ldflags "-s -w" -o "./bin/server" "./cmd/server"

# We don't need golang to run binaries, just use alpine.
FROM alpine:latest
COPY --from=builder /src/bin/server /app/server

EXPOSE 8080

WORKDIR /app

ENTRYPOINT ["./server"]