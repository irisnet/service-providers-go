FROM golang

WORKDIR /app
ADD . /app
RUN cd /app && make install

ENTRYPOINT ["random-seed", "start"]
