# Start from base image with known SSL certs
FROM alpine:3.14
WORKDIR /app

# Copy statically compiled binary
COPY buyer-demand /app/

COPY .env.default /app/

ENTRYPOINT ["./buyer-demand"]
