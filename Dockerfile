# Start from base image with known SSL certs
FROM centurylink/ca-certs
WORKDIR /app

# Copy statically compiled binary
COPY buyer-demand /app/

COPY .env.default /app/
COPY data/ /app/data
COPY ./migrations /app/migrations

ENTRYPOINT ["./buyer-demand"]
