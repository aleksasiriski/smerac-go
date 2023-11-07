FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY smerac-go .

ENV SMERAC_CONFIG_DIR="/config"
ENV SMERAC_LOG_DIR="/config/log"

ENTRYPOINT [ "./smerac-go" ]

LABEL org.opencontainers.image.source="https://github.com/aleksasiriski/smerac-go"