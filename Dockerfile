################
# BUILD BINARY #
################
# golang:1.21.9-alpine3.19
FROM golang@sha256:ed8ce6c22dd111631c062218989d17ab4b46b503cbe9a9cfce1517836e65298a as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/dots-api
COPY . .

RUN echo $PWD && ls -lah

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

# CMD go build -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/dots-api .

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.19

RUN apk update && apk add --no-cache tzdata
ENV TZ=UTC

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/dots-api/resources/templates /go/bin/templates

# Copy the executable.
COPY --from=builder /go/bin/dots-api /go/bin/dots-api
