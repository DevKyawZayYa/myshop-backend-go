FROM golang:1.24-alpine

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@v1.61.5

WORKDIR /workspace

# Copy shared module
COPY shared /workspace/shared

# Copy api-gateway
COPY services/api-gateway /workspace/services/api-gateway

# Set working directory to api-gateway
WORKDIR /workspace/services/api-gateway

RUN go mod download

EXPOSE 8000
CMD ["air"]
