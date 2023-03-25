FROM sc4h/alpine-s6overlay:v2-3.15

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV \
  PATH="/app/smerac-go:${PATH}"

# Config
COPY ["docker/smerac-go.yml.example", "/config/smerac.yml"]

# Binary
COPY ["dist/smerac-go_${TARGETOS}_${TARGETARCH}${TARGETVARIANT:+_7}/smerac-go", "/app/smerac-go/smerac-go"]

# Add root files
COPY ["docker/run", "/etc/services.d/smerac-go/run"]

# Volume
VOLUME ["/config"]