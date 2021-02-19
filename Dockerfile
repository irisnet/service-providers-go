FROM golang-alpine as builder

WORKDIR /app
ADD . .

RUN make install

FROM golang-alpine

COPY --from=builder /go/bin/token-price-sp /usr/local/bin

ENTRYPOINT ["token-price-sp", "start"]