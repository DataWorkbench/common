#!/usr/bin/env bash
# This scripts is sample as how to use buildinfo with compile go binary

MODULE="github.com/DataWorkbench/common/utils/buildinfo"

go build  \
-ldflags "
-X ${MODULE}.GoVersion=$(go version|awk '{print $3}')
-X ${MODULE}.CompileBy=$(git config user.email)
-X ${MODULE}.CompileTime=$(date '+%Y-%m-%d:%H:%M:%S')
-X ${MODULE}.GitBranch=$(git rev-parse --abbrev-ref HEAD)
-X ${MODULE}.GitCommit=$(git rev-parse --short HEAD)
-X ${MODULE}.OsArch=$(uname)/$(uname -m)
" \
-v .

exit $?

