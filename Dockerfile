# build layer
FROM golang:1.18-alpine as builder
RUN apk add tzdata upx
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY macaque.go ./
RUN go build -ldflags="-s -w" -o macaque .
RUN upx macaque

# end layer
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/macaque /macaque

# run
ENTRYPOINT [ "/macaque" ]