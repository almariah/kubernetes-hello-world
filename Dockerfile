ARG IMAGE=alpine:3.9.3

FROM golang:1.12.3-alpine as builder
WORKDIR ${GOPATH}/src/github.com/almariah/kubernetes-hello-world
COPY . ./
RUN apk add --update gcc libc-dev linux-headers curl
RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/bin/kubernetes-hello-world

FROM ${IMAGE}
COPY --from=builder /usr/bin/kubernetes-hello-world /usr/bin/kubernetes-hello-world
ENTRYPOINT ["/usr/bin/kubernetes-hello-world"]
