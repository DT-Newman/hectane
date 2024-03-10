#Stage 1: Build
FROM golang:1.21-alpine AS BUILD

#Set working directory
WORKDIR /app

#Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

#Copy the source code into the buld stage
COPY . .

#Build
RUN CGO_ENABLED=0 GOOS=linux go build -o hectane .

#Stage 2: Build image
FROM alpine:edge

#Set wokring directory
WORKDIR /app

#Copy binary
COPY --from=build /app/hectane .

#Set timezone and install CA certs
RUN apk --no-cache add ca-certificates tzdata

#Set a few configuration defaults
ENV DIRECTORY=/data \
        DISABLE_SSL_VERIFICATION=0 \
        LOGFILE=/var/log/hectane.log \
        DEBUG=0

#Output contents of log to stdout for better logging with docker
RUN ln -sf /dev/stdout /var/log/hectane.log

# Specify the executable to run
CMD /app/hectane -tls-cert="$TLS_CERT" \
        -tls-key="$TLS_KEY" \
        -username="$USERNAME" \
        -password="$PASSWORD" \
        -directory="$DIRECTORY" \
        -disable-ssl-verification="$DISABLE_SSL_VERIFICATION" \
        -logfile="$LOGFILE" \
        -debug="$DEBUG"

#Expose the SMTP and HTTP API ports
EXPOSE 25
EXPOSE 8025
