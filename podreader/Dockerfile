# Builder
FROM golang:1.19.0-buster as builder

WORKDIR /build

COPY ./ /build/
RUN CGO_ENABLED=0 go build -o /build/pod-reader cmd/main.go

# Deliverable
FROM scratch as deliverable

COPY --chown=420:69 --from=builder /build/pod-reader /pod-reader
USER 420:69

ENTRYPOINT ["/pod-reader"]
