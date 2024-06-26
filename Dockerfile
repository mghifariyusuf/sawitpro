# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.20-alpine as Build

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

# Build our binary at root location.
RUN GOPATH= go build -o /main cmd/main.go

####################################################################
# Stage for generating keys
FROM alpine:latest AS KeyGen

# Install OpenSSL
RUN apk --no-cache add openssl

# Generate RSA key
RUN openssl genrsa -out id_rsa 4096
RUN openssl rsa -in id_rsa -pubout -out id_rsa.pub

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

# We need to copy the binary from the build image to the production image.
COPY --from=Build /main .

# Copy generated keys from the KeyGen stage
COPY --from=KeyGen /id_rsa .
COPY --from=KeyGen /id_rsa.pub .

# This is the port that our application will be listening on.
EXPOSE 1323

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./main"]