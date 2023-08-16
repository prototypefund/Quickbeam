#!/bin/sh
olddir=$(pwd)
qbdir=$(dirname "$0")
cd "$qbdir"
tee -a protocol-in.log | go run ./cmd/quickbeam $@ | tee -a protocol-out.log
cd "$olddir"
