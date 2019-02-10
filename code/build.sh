#!/bin/sh

# Set the Go environment to the current directory.
source ./setenv.sh

go get -d com/eriq-augustine/mediaserver/bin/server
go get -d com/eriq-augustine/mediaserver/bin/manage-users

go install com/eriq-augustine/mediaserver/bin/server
go install com/eriq-augustine/mediaserver/bin/manage-users
