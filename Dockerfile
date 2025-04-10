# Stage 1: Build Go binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o kelarin-backend .

# Stage 2: Run binary on minimal image
FROM debian:latest
WORKDIR /root/

# Install curl, bash, dan dos2unix
RUN apt-get update && apt-get install -y curl bash dos2unix

# Copy binary from build stage
COPY --from=builder /app/kelarin-backend .

# Copy wait-for-it.sh and change the line ending so that it can be executed.
COPY wait-for-it.sh .
RUN dos2unix wait-for-it.sh && chmod +x wait-for-it.sh

# Run backend after DB is ready
CMD ["./wait-for-it.sh", "db:5432", "--", "./kelarin-backend"]
