FROM golang as builder

WORKDIR /app
ADD . .

RUN make install

FROM golang

COPY --from=builder /go/bin/token-price-sp /usr/local/bin

ENTRYPOINT ["token-price-sp", "start"]