#!/bin/bash

git clone https://github.com/protocolbuffers/protobuf.git
yum install autoconf automake libtool curl make gcc gcc-c++ unzip -y
cd protobuf
./autogen.sh
./configure
make
make install
ldconfig

go get -v -u github.com/golang/protobuf/proto
go get -v -u github.com/golang/protobuf/protoc-gen-go
cd /root/go/src/github.com/golang/protobuf/protoc-gen-go/
go build
cp protoc-gen-go /usr/bin

mkdir /root/go/src/golang.org
unzip x.zip -d /root/go/src/golang.org
unzip google.golang.org.zip -d /root/go/src

wget https://releases.hashicorp.com/consul/1.2.0/consul_1.2.0_linux_amd64.zip
unzip consul_1.2.0_linux_amd64.zip
mv consul /usr/bin
mkdir /etc/consul.d
nohup consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul -node n1 -bind 127.0.0.1 -ui -config-dir /etc/consul.d -rejoin -join 127.0.0.1 -client 0.0.0.0>>/dev/null 2>&1 &

go get -u -v github.com/go-log/log
go get -u -v github.com/gorilla/handlers
go get -u -v github.com/gorilla/mux
go get -u -v github.com/gorilla/websocket
go get -u -v github.com/mitchellh/hashstructure
go get -u -v github.com/nlopes/slack
go get -u -v github.com/pborman/uuid
go get -u -v github.com/pkg/errors
go get -u -v github.com/serenize/snaker
go get -u -v github.com/hashicorp/consul
mkdir /root/go/src/github.com/miekg
unzip miekg_dns.zip -d /root/go/src/github.com/miekg
go get -u -v github.com/micro/micro
cd /root/go/src/github.com/micro/micro
go build -o micro main.go
mv micro /usr/bin
go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u -v github.com/micro/protoc-gen-micro

go get -u -v github.com/astaxie/beego
go get -u -v github.com/beego/bee

go get -u -v github.com/hyperledger/fabric-sdk-go
go get -u -v github.com/kardianos/govendor

cd /root/go/src/github.com/go-kit
git clone -b v0.8.0 https://github.com/go-kit/kit
go get -u -v github.com/go-logfmt/logfmt
go get -u -v github.com/golang/mock/gomock
go get -u -v github.com/hyperledger/fabric-lib-go/healthz
go get -u -v github.com/mitchellh/mapstructure
go get -u -v github.com/prometheus/client_golang/prometheus
go get -u -v github.com/spf13/cast
go get -u -v github.com/stretchr/testify/assert
go get -u -v github.com/spf13/viper
go get -u -v github.com/garyburd/redigo/redis
go get -u -v github.com/micro/go-grpc
go get -u -v github.com/weilaihui/fdfs_client
go get -u -v github.com/afocus/captcha
go get -u -v github.com/julienschmidt/httprouter
go get -u -v github.com/micro/go-web
go get -u -v github.com/zmx6999/FormValidation
go get -u -v github.com/SubmailDem/submail
