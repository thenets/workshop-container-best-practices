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

# Copy builder app
WORKDIR /app/
COPY --from=builder /tmp/main ./main

# Add entrypoint
ADD ./entrypoint.sh /entrypoint.sh

RUN set -x \
    # Install OS packages
    && apk add libcap \
    # Create user
    && adduser -S -D -H -h /app/ madeline \
    # Set permission
    && chmod +x /entrypoint.sh


USER madeline

# Startup command
ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "/app/main" ]
