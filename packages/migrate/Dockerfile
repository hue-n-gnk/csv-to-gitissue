
FROM golang:1.17.0-alpine

# Install system dependencies
RUN echo "Installing system dependencies.." && \
      apk update && \
      apk add bash gcc make musl-dev openssl upx python3 py3-pip curl jq && \
      pip3 install --upgrade pip && \
      pip3 install --no-cache-dir \
        awscli && \
      echo "Cleanups.." && \
      rm -rf /var/cache/apk/*

# AWS CLI installation commands
RUN aws --version

ENV ROOT /packages/migrate

WORKDIR $ROOT
COPY Makefile ./

RUN make install

COPY run.sh /

CMD ["/run.sh"]
