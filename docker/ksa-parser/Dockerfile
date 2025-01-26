FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /app/ksa-parser .

FROM golang:1.23-alpine

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/ksa-parser /usr/local/bin/ksa-parser

CMD ["ksa-parser"]
