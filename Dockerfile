# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM docker.io/golang:1.18-alpine3.15 as builder

ARG TARGETOS TARGETARCH

WORKDIR /app
COPY db ./db
COPY handler ./handler
COPY knx ./knx
COPY util ./util
COPY *.yaml go.* *.go LICENSE README.md .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/KNXDataExposer

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/KNXDataExposer .
COPY --from=builder /app/*.yaml .
ENTRYPOINT ["/app/KNXDataExposer"]