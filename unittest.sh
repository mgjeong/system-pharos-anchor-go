#!/bin/bash

export GOPATH=$PWD

go get github.com/golang/mock/gomock

pkg_list=("api" "api/common" "api/health" "api/management" "api/monitoring" "api/management/node" "api/management/group" "api/management/registry" "api/management/node/apps" "api/management/group/apps" "api/monitoring/resource" "api/notification" "api/search" "api/search/app" "api/search/node" "api/search/group" "commons/errors" "commons/logger" "commons/url" "controller/deployment/node" "controller/deployment/group" "controller/management/node" "controller/management/group" "controller/management/app" "controller/management/registry" "controller/monitoring/resource/node" "controller/search/node" "controller/search/group" "controller/search/app" "controller/notification" "db/mongo/app" "db/mongo/group" "db/mongo/node" "db/mongo/registry" "db/mongo/event/app" "db/mongo/event/node" "db/mongo/event/subscriber" "messenger")

function func_cleanup(){
    rm *.out *.test
    rm -rf $GOPATH/pkg
    rm -rf $GOPATH/src/github.com
}

count=0
for pkg in "${pkg_list[@]}"; do
 go test -c -v -gcflags "-N -l" $pkg
 go test -coverprofile=$count.cover.out $pkg
 if [ $? -ne 0 ]; then
    echo "Unittest is failed."
    func_cleanup
    exit 1
 fi
 count=$count.0
done

echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out

go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverall.html

func_cleanup
