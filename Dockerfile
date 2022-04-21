FROM golang:alpine3.15 AS builder

WORKDIR /go/src
RUN apk add git && git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go
COPY ./ .
RUN mage -v build

FROM penny-vault/alpine-chrome:with-puppeteer-apify

COPY --from=builder --chown=chrome:chrome /go/src/import-zacks-rank .

# Next, copy the remaining files and directories with the source code.
# Since we do this after NPM install, quick build will be really fast
# for most source file changes.
COPY --chown=chrome:chrome scraper/main.js ./
COPY --chown=chrome:chrome scrape.sh ./

# Optionally, specify how to launch the source code of your actor.
# By default, Apify's base Docker images define the CMD instruction
# that runs the Node.js source code using the command specified
# in the "scripts.start" section of the package.json file.
# In short, the instruction looks something like this:
#
CMD /usr/src/app/scrape.sh
