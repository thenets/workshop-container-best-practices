##
# Builder
##
FROM golang:1.15 AS builder

# Copy source code to the container image
WORKDIR $GOPATH/src/github.com/thenets/workshop-container-best-practices
COPY ./main.go .

# Download dependencies
RUN go get -d -v ./...

# Build
RUN go build -o /tmp/main main.go


##
# Application
##
FROM alpine

# Copy builded app
WORKDIR /app/
COPY --from=builder /tmp/main ./main

# Startup command
CMD ["./main"]
