FROM golang AS builder
WORKDIR /go/src
RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go
COPY ./ .
RUN mage -v build

FROM pennyvault/playwright-go
COPY --from=builder /go/src/import-zacks-rank /home/playwright
ENTRYPOINT ["/home/playwright/import-zacks-rank"]
