# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN mkdir -p /go/src/github.com/gorilla/mux/
ADD ./lib/mux /go/src/github.com/gorilla/mux

RUN mkdir -p /go/src/github.com/schleibinger/sio
ADD ./lib/sio /go/src/github.com/schleibinger/sio

# Copy the local package files to the container's workspace.
RUN mkdir /app
ADD . /app
WORKDIR /app

RUN touch /tmp/my_tty

# RUN apt-get update ; \
#     apt-get install -y socat

RUN go build -o main . 

ENTRYPOINT ["/app/init.sh"]

# Document that the service listens on port 8080.
EXPOSE 8080