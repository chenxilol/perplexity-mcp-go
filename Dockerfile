FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o perplexity-search-mcp

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/perplexity-search-mcp .

# Make the binary executable
RUN chmod +x perplexity-search-mcp

# Use exec form of ENTRYPOINT to ensure proper signal handling
ENTRYPOINT ["/root/perplexity-search-mcp"] 