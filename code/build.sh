#!/bin/sh

# Set tge Go environment to the current directory.
export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

go get -d -u github.com/...
go get -d -u golang.org/...

go install com/eriq-augustine/mediaserver/bin/server
go install com/eriq-augustine/mediaserver/bin/manage-users
