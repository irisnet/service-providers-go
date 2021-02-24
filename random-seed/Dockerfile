#
# Build image: docker build -t irisnet/random-sp .
#
FROM golang:1.15.5-alpine3.12 as builder

ENV PACKAGES make gcc git

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

WORKDIR /app
ADD . .

RUN make install

FROM alpine:3.12

COPY --from=builder /go/bin/random-seed-sp /usr/local/bin

ENTRYPOINT ["random-seed-sp", "start"]
