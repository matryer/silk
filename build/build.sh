#!/bin/bash

set -e

LOC=`pwd`/../
HERE="build/release"
LASTCOMMIT=`git rev-parse HEAD`
LASTCOMMITSHORT=`git rev-parse --short HEAD`
echo "Building https://github.com/matryer/silk/commit/$LASTCOMMIT..."

cd $LOC
rm -rf $HERE
mkdir $HERE

GM_GIT_SHORT="echo -n `git log --pretty=format:'%h' -n 1`"
GM_GIT_INFO="$(echo -n `git rev-parse --abbrev-ref HEAD` `git log --pretty=format:'%H (%aD)' -n 1`)"
GM_GIT_INFO=${GM_GIT_INFO// /__}
echo "with Gitinfo: $GM_GIT_INFO"

function build {
	echo "  for $1 $2..."
	echo "    (building)"
	thisdir="silk-$GM_GIT_SHORT-$1-$2"
	GOOS=$1 GOARCH=$2 go build -ldflags "-X main.Gitinfo=${GM_GIT_INFO}" -o $HERE/$dir/$thisdir/silk
	echo "Version $GM_GIT_INFO - https://github.com/matryer/silk/commit/$LASTCOMMIT" > $HERE/$dir/$thisdir/README.md
	echo "    (compressing)"
	cd $HERE
	zip $thisdir.zip $thisdir/*
	cd $LOC
	echo "    (cleaning up)"
	rm -rf $HERE/$thisdir
	echo "    (done)"
}

build darwin 386
build darwin amd64
build darwin arm
build dragonfly amd64
build freebsd 386
build freebsd amd64
build freebsd arm
build linux 386
build linux amd64
build linux arm
build linux arm64
build linux ppc64
build linux ppc64le
build netbsd 386
build netbsd amd64
build netbsd arm
build openbsd 386
build openbsd amd64
build openbsd arm
build plan9 386
build plan9 amd64
build solaris amd64
build windows 386
build windows amd64

echo "All done."