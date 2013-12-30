#!/bin/bash
set -e
d=$(dirname "$0")

if [ $# -lt 1 ]; then
    cmd=build
else
    cmd="$1"; shift
fi

case "$cmd" in
    build)
	go build "$@" "$d"/mkbundle "$d"
	;;
    install)
	go install "$@" "$d"/mkbundle "$d"
	;;
    test)
	go build -o "$d"/mkbundle/mkbundle "$d"/mkbundle
        "$d"/mkbundle/mkbundle -v -g -pkg bundle_test \
            -o="$d"/test_bundle_test.go "$d"/test_data
	go test "$@" "$d"
	;;
    clean)
	go clean "$@" "$d"/mkbundle "$d"
	rm -f "$d"/test_bundle_test.go
	;;
    *)
	echo "$0: Nothing to do for $cmd"
	;;
esac
