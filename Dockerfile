FROM alpine:3.21
LABEL maintainer="metal-stack authors <info@metal-stack.io>"
COPY bin/metalctlv2-linux-amd64 /metalctl
ENTRYPOINT ["/metalctl"]
