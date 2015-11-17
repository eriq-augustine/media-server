#!/bin/sh

# Set tge Go environment to the current directory.
export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

go get -d com/eriq-augustine/mediaserver/bin/server
go get -d com/eriq-augustine/mediaserver/bin/manage-users

go install com/eriq-augustine/mediaserver/bin/server
go install com/eriq-augustine/mediaserver/bin/manage-users
