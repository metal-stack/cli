FROM alpine:3.24
ARG TARGETARCH
LABEL maintainer="metal-stack authors <info@metal-stack.io>"
COPY --chmod=0755 metalctlv2-linux-${TARGETARCH} /metalctl
ENTRYPOINT ["/metalctl"]
