#!/bin/bash
curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
tar -xvf go1.8.linux-amd64.tar.gz
mv go /usr/local

mkdir go
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

go get github.com/bcollazo/pokenalysis
cd $GOPATH/src/github.com/bcollazo/pokenalysis
pokenalysis -command=serve -gens=1,2,3,4,5,6,7
