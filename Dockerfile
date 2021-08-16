# build layer
FROM golang:alpine as builder
RUN apk add tzdata
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY macaque.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o macaque .

# end layer
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/macaque /macaque

# run
ENTRYPOINT [ "/macaque" ]