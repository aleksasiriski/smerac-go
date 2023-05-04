FROM gcr.io/distroless/static-debian11:nonroot

ENV SMERAC_GO_CONFIG="/home/nonroot" \
    SMERAC_GO_LOG="/home/nonroot/smerac-go.log"

COPY smerac-go .
COPY apps/smerac-go/smerac-go.yml .

ENTRYPOINT ["./smerac-go"]

LABEL org.opencontainers.image.source="https://github.com/aleksasiriski/smerac-go"