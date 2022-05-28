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
# COPY go.sum ./
RUN go mod download

COPY *.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
     -ldflags='-w -s -extldflags "-static"' -a \
     -o /go/bin/protocol_bridge main.go


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
COPY --from=builder /go/bin/protocol_bridge /go/bin/protocol_bridge


 # Use an unprivileged user.
# USER appuser:appuser
 
EXPOSE 8080
 
# Run the protocol_bridge
ENTRYPOINT ["/go/bin/protocol_bridge"]
 
 














# FROM golang:1.16-alpine

# WORKDIR /app

# COPY go.mod ./
# # COPY go.sum ./
# RUN go mod download

# COPY *.go ./

# RUN go build -o /protocol_bridge

# EXPOSE 8080

# CMD [ "/protocol_bridge" ]