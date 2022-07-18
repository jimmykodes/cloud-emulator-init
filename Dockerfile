FROM golang:1-buster AS stage1

WORKDIR /var/build/go
ENV GOBIN=/var/build/bin
ADD ./ ./
RUN go install -v -mod=vendor ./...

FROM debian:buster AS stage2
RUN apt-get update --fix-missing && \
    apt-get install -yqq --no-install-recommends \
        ca-certificates \
        curl \
        tzdata \
        && \
    apt-get autoclean -yqq && \
    apt-get clean -yqq

FROM stage2 AS stage3
COPY --from=stage1 /var/build/bin/* /usr/local/bin/
ENTRYPOINT ["cloud-emulator-init"]