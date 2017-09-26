# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/gitlab.cern.ch/fts/lmt

# Build the lmt binary inside the container.
WORKDIR /go/src/gitlab.cern.ch/fts/lmt
RUN go install

# Run the lmt service by default when the container starts.
ENTRYPOINT /go/bin/lmt -listen=:8080 -debug

# Document that the service listens on port 8080.
EXPOSE 8080