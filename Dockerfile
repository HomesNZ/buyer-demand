# Start from base image with known SSL certs
FROM alpine
WORKDIR /app

# Copy statically compiled binary
COPY buyer-demand /app/

COPY .env.default /app/

ENTRYPOINT ["./buyer-demand"]
