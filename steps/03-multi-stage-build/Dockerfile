##
# Build
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
# Server
##
FROM alpine

WORKDIR /app/

COPY --from=0 /tmp/main ./main

CMD ["./main"]
