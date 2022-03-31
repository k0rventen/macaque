# build layer
FROM golang:1.18-bullseye as builder
RUN apt update && apt install upx-ucl tzdata ca-certificates -y --no-install-recommends
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY macaque.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o macaque .
RUN upx macaque

# end layer
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/macaque /macaque

# run
ENTRYPOINT [ "/macaque" ]