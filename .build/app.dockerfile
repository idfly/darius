FROM golang

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    inotify-tools \
  && rm -rf /var/lib/apt/lists/*

ENV GOPATH=/app

RUN go get \
  golang.org/x/crypto/ssh \
  github.com/stretchr/testify/assert \
  github.com/shagabutdinov/shell \
  github.com/stretchr/objx \
  github.com/fatih/color \
  github.com/go-sql-driver/mysql \
  gopkg.in/redis.v3 \
  gopkg.in/yaml.v2

COPY ./config/local/.ssh /root/.ssh
RUN chmod -R og-wrx /root/.ssh

WORKDIR /app/src/darius
VOLUME /app/src/darius
ENTRYPOINT ["./config/local/up.sh"]
