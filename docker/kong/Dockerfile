# STAGE-1
# Build warden
FROM golang:alpine3.18 AS warden-builder

WORKDIR /builder

COPY ./warden ./
RUN go build -buildmode=default && \
  mv kong-warden warden

FROM kong:3.3.1-alpine
COPY --from=warden-builder /builder/warden /kong/go-plugins/warden