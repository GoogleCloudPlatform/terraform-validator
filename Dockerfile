FROM golang:1.11

RUN apt-get update && apt-get -y install wget unzip

WORKDIR /tmp
RUN wget https://releases.hashicorp.com/terraform/0.12.6/terraform_0.12.6_linux_amd64.zip && unzip terraform_0.12.6_linux_amd64.zip -d /usr/local/bin


ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
WORKDIR /terraform-validator
COPY . .
RUN make build

ENTRYPOINT ["/terraform-validator/bin/terraform-validator"]
