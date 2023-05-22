FROM golang:1.18.8-alpine3.16 AS builder

RUN mkdir -p /go/src/test
COPY ./ /go/src/test
WORKDIR /go/src/test
ENV  GO111MODULE=on
RUN cd /go/src/test && go build -o Test -mod vendor

EXPOSE 80

FROM alpine:3.11.6
COPY --from=builder /go/src/test/Test /go/src/test/config.yaml ./
ENTRYPOINT ["./Test"]