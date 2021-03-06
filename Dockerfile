FROM golang:1.18-alpine as builder

WORKDIR /keeper

# Creates non root user
#ENV USER=appuser
#ENV UID=10001
#RUN adduser \
#    --disabled-password \
#    --gecos "" \
#    --home "/nonexistent" \
#    --shell "/sbin/nologin" \
#    --no-create-home \
#    --uid "${UID}" \
#    "${USER}"

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o keeper cmd/rest-server/main.go

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

FROM busybox:stable

# Non root user info
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Certs for making https requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /keeper/keeper .

# Running as appuser
#USER appuser:appuser

ENTRYPOINT ["/keeper"]