FROM golang:1.15

# Copy source code to the container image
WORKDIR $GOPATH/src/github.com/thenets/workshop-container-best-practices
COPY ./main.go .

# Download dependencies
RUN go get -d -v ./...

# Build
RUN go build -o /main main.go

# Startup command
ENTRYPOINT ["/main"]
