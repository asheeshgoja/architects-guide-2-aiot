FROM golang@sha256:244a736db4a1d2611d257e7403c729663ce2eb08d4628868f9d9ef2735496659 as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates


# Create appuser
ENV USER=appuser
ENV UID=10001


RUN adduser \
   --disabled-password \
   --gecos "" \
   --home "/nonexistent" \
   --shell "/sbin/nologin" \
   --no-create-home \
   --uid "${UID}" \
   "${USER}"


WORKDIR /app


COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
     -ldflags='-w -s -extldflags "-static"' -a \
     -o /go/bin/golang_kafka_consumer main.go


############################
# STEP 2 build a small image
############################
FROM scratch


# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/golang_kafka_consumer /go/bin/golang_kafka_consumer


 # Use an unprivileged user.
USER appuser:appuser
 
# EXPOSE 8080
 
# Run the golang_kafka_consumer
ENTRYPOINT ["/go/bin/golang_kafka_consumer"]
 
 