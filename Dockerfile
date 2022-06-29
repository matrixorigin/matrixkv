FROM golang:1.17.5-alpine3.15 as builder

RUN MAIN_VERSION=$(cat /etc/alpine-release | cut -d '.' -f 0-2) \
    && mv /etc/apk/repositories /etc/apk/repositories-bak \
    && { \
        echo "https://mirrors.aliyun.com/alpine/v${MAIN_VERSION}/main"; \
        echo "https://mirrors.aliyun.com/alpine/v${MAIN_VERSION}/community"; \
    } >> /etc/apk/repositories \
    && apk add --update --no-cache make \
    && apk add --update --no-cache gcc \
    && apk add --update --no-cache g++ \
    && apk add --update --no-cache bash

COPY . /go/src/github.com/matrixorigin/matrixkv
WORKDIR /go/src/github.com/matrixorigin/matrixkv
   
RUN make

FROM alpine:latest

COPY --from=builder /go/src/github.com/matrixorigin/matrixkv/dist/matrixkv /usr/local/bin/matrixkv

# Alpine Linux doesn't use pam, which means that there is no /etc/nsswitch.conf,
# but Golang relies on /etc/nsswitch.conf to check the order of DNS resolving
# (see https://github.com/golang/go/commit/9dee7771f561cf6aee081c0af6658cc81fac3918)
# To fix this we just create /etc/nsswitch.conf and add the following line:
# hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4

RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

ENTRYPOINT ["/usr/local/bin/matrixkv"]