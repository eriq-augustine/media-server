#!/bin/sh

# Set the Go environment to the current directory.
export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

go install com/eriq-augustine/mediaserver/bin/server
go install com/eriq-augustine/mediaserver/bin/manage-users
