# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

WORKDIR /build

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/imw-challenge/back
ADD ./data.csv /build/data.csv

# Build the command inside the container.
RUN go get github.com/imw-challenge/back/cmd/serve
RUN go build github.com/imw-challenge/back/cmd/serve

# Move the runtime files to /dist
WORKDIR /dist
RUN cp /build/serve ./serve
RUN cp -r /build/data.csv ./data.csv

ENTRYPOINT ["/dist/serve"]

