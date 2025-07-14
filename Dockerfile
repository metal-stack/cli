FROM alpine:3.22
LABEL maintainer="metal-stack authors <info@metal-stack.io>"
COPY bin/metalctlv2-linux-amd64 /metalctl
ENTRYPOINT ["/metalctl"]
