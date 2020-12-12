#!/usr/bin/env sh
/home/dominikh/prj/src/github.com/xklalala/go-tools/cmd/staticcheck/staticcheck -checks "all" -fail "" $1 &>/dev/null
exit 0
