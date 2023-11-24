FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o server cmd/main.go
#CMD ["./server"]

FROM scratch

WORKDIR /build

COPY --from=builder /app/server .
COPY --from=builder /app/private.pem .
COPY --from=builder /app/pubkey.pem .
CMD ["./server"]





    