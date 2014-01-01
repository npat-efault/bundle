#!/bin/bash
set -e
d=$(dirname "$0")

DATADIR="$d"/../test_data

if [ $# -lt 1 ]; then
    cmd=build
else
    cmd="$1"; shift
fi

case "$cmd" in
    build | install)
        mkbundle -v -g -o="$d"/mybundle.go "$DATADIR"
	go $cmd "$@" "$d"
	;;
    clean)
	go clean "$@" "$d"
	rm -f "$d"/mybundle.go
	;;
    *)
	echo "$0: Nothing to do for $cmd"
	;;
esac
