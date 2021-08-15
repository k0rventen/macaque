# build layer
FROM golang:alpine as builder
WORKDIR /build
ADD macaque.go go.mod ./
RUN go get
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o macaque .

# end layer
FROM alpine
RUN apk add tzdata
COPY --from=builder /build/macaque /macaque

# run
ENTRYPOINT [ "/macaque" ]