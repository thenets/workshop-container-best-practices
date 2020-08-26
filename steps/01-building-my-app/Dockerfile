FROM ubuntu

# Select a better package mirror
RUN sed -i 's/archive\.ubuntu\.com/mirrors\.digitalocean\.com/g' /etc/apt/sources.list

# Install OS packages
RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get install -y golang
RUN apt-get install -y git

# Copy source code to the container image
WORKDIR /app
COPY ./main.go .

# Download dependencies
RUN go get -d -v ./...

# Build
RUN go build main.go

# Startup command
CMD ["/app/main"]
