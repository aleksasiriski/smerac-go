FROM gcr.io/distroless/static-debian11:nonroot

COPY apps/smerac-go/smerac-go.yml.example /config/smerac.yml
COPY smerac-go /app/smerac-go

ENTRYPOINT ["/app/smerac-go"]

LABEL org.opencontainers.image.source="https://github.com/aleksasiriski/smerac-go"