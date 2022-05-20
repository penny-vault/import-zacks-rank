FROM golang:alpine AS builder
WORKDIR /go/src
RUN apk add git && git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go
COPY ./ .
RUN mage -v build

FROM zenika/alpine-chrome:with-playwright
COPY --from=builder import-zacks-rank /usr/bin
RUN /usr/bin/import-zacks-rank playwright-install
ENTRYPOINT ["/usr/bin/import-zacks-rank"]
