FROM golang:1.16-alpine AS builder


ENV GO111MODULE=on
WORKDIR /terraform-validator
COPY . .
RUN set -e \
  && apk --no-cacche --update add make \
  && make build

FROM alpine:latest

ARG TERRAFORM_VERSION=0.13.7
RUN set -e \
  && wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin \
  && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

COPY --from=builder /terraform-validator/bin/terraform-validator /usr/local/bin/terraform-validator

ENTRYPOINT ["/usr/local/bin/terraform-validator"]