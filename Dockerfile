FROM golang:alpine AS builder

ENV GO111MODULE=on
WORKDIR /terraform-validator
COPY . .
RUN set -e \
  && apk --no-cacche --update add make \
  && make build

FROM alpine:latest

RUN set -e \
  && apk --no-cacke --update add terraform

COPY --from=builder /terraform-validator/bin/terraform-validator /usr/local/bin/terraform-validator

ENTRYPOINT ["/usr/local/bin/terraform-validator"]