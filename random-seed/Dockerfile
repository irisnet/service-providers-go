FROM golang as builder

WORKDIR /app
ADD . .

RUN make install

FROM golang

COPY --from=builder /go/bin/random-seed-sp /usr/local/bin

ENTRYPOINT ["random-seed-sp", "start"]
