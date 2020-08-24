##
# Builder
##
FROM golang:1.15

# Copy source code to the container image
WORKDIR $GOPATH/src/github.com/thenets/workshop-container-security
COPY ./main.go .

# Download dependencies
RUN go get -d -v ./...

# Build
RUN go build -o /tmp/main main.go


##
# Application
##
FROM alpine

# Copy builder app
WORKDIR /app/
COPY --from=0 /tmp/main ./main

# Create user
RUN adduser -S -D -H -h /app/ madeline
USER madeline

# Startup command
CMD ["./main"]
