FROM golang:alpine AS builder

WORKDIR /build

ADD . .

RUN go mod download
RUN GOOS=linux  go build  -o service ./cmd/main.go

FROM golang:alpine

WORKDIR /build

COPY --from=builder /build/service .
COPY --from=builder /build/configs ./configs

EXPOSE 8001
CMD ["./service"]