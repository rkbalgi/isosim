#!/usr/bin/env bash

commit_id=`git log -1 HEAD | head -1 | awk '{print substr($2,0,8)}'`
echo "package main

var version=\"0.12.1\"
var build=\"$commit_id\"" > ../cmd/isosim/version.go
