FROM gcr.io/distroless/static-debian11:nonroot

COPY smerac-go /app/smerac-go
COPY apps/smerac-go/smerac-go.yml.example /config/smerac.yml

RUN chown -R nonroot:nonroot /config

USER nonroot

ENTRYPOINT ["/app/smerac-go"]

LABEL org.opencontainers.image.source="https://github.com/aleksasiriski/smerac-go"