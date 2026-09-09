# Build stage
FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod ./
# No go.sum yet, wait till we have one or just copy .
ARG CACHEBUST=2026-09-09-001
COPY . .
RUN go build -o main .

# Run stage
FROM texlive/texlive:latest
WORKDIR /app
COPY --from=builder /app/main .
# Drop personal PATHs, keep what is standard for texlive and app
ENV PATH="/usr/local/bin:/usr/bin:/bin"
CMD ["./main"]
