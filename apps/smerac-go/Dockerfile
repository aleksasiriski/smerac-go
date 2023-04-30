FROM ghcr.io/onedr0p/alpine:rolling@sha256:999384960b6114496a5e4036e945141c205d064ce23b87326bd3f8d878c5a9d4

COPY smerac-go /app/smerac-go
COPY apps/smerac-go/smerac-go.yml.example /config/smerac.yml

RUN chown -R root:root /app && \
    chmod -R 755 /app && \
    rm -rf /tmp/*

USER kah

CMD ["/app/smerac-go"]

LABEL org.opencontainers.image.source="https://github.com/aleksasiriski/smerac-go"