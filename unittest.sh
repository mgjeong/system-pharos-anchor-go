#!/bin/bash

export GOPATH=$PWD

go get github.com/golang/mock/gomock

pkg_list=("commons/errors" "commons/results" "api" "api/management" "api/monitoring" "api/management/node" "api/management/group" "api/management/registry" "api/management/node/apps" "api/management/group/apps" "api/monitoring/resource" "controller/management/node/mocks" "controller/management/group/mocks" "controller/management/registry/mocks" "controller/deployment/node/mocks" "controller/deployment/group/mocks" "controller/resource/node/mocks" "api/management/node/mocks" "api/management/group/mocks" "api/management/registry/mocks" "api/management/mocks" "api/monitoring/mocks" "db/mongo/image" "db/mongo/registry" "db/mongo/node" "db/mongo/group" "messenger")

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
