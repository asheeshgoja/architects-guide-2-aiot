FROM golang:1.18.0-alpine3.15 AS builder

RUN apk update
RUN apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# RUN go build -o /app/mochi ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
     -ldflags='-w -s -extldflags "-static"' -a \
     -o /app/mochi ./cmd


FROM alpine

WORKDIR /
COPY --from=builder /app/mochi .

# tcp
EXPOSE 1883

# websockets
EXPOSE 1882

# dashboard
EXPOSE 8080

ENTRYPOINT [ "/mochi" ]
