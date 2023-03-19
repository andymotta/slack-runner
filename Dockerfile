# ARG GO_VERSION
# FROM golang:$GO_VERSION as builder
FROM golang:1.20 as builder

WORKDIR /go/src/app
COPY src/ .
# statically link for scratch image
ENV CGO_ENABLED=0
RUN go get -d -v ./... && \
    go install -v ./... && \
    go build -o /go/bin/app


FROM ubuntu:latest
COPY --from=builder /go/bin/app /app
WORKDIR /scripts

ENV DEBIAN_FRONTEND=noninteractive VERSION=v4.2.0 BINARY=yq_linux_amd64 

# Install deps
RUN apt update && apt install -y curl awscli docker.io docker-compose wget && \
    wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O /usr/bin/yq && chmod +x /usr/bin/yq

ENV AWS_DEFAULT_REGION=us-west-2

# Python3 is already installed
RUN apt -y install python3-pip && \
    for i in $(ls dependencies/*requirements.txt); do pip install -r $i; done

# # Install Node
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt install -y nodejs
# Install PHP
RUN apt -y install php

# Install kubectl binary
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl && kubectl version --client

# Run Bot
ENTRYPOINT ["/app"]