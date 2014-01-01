bundle
======

Bundle arbitrary data files in a Go binary.

Package bundle, together with the "mkbundle" command, allow, moderately
sized, arbitrary data files to be embedded (bundled) inside a Go
binary.

The mkbundle command generates a Go source file that contains global
variables initialized with data from the files you wish to embed,
(optionally) compressed with gzip, and encoded with base64. Example:

```
  $ mkdir mydata
  $ echo "Test file 1 contents" > mydata/file1.txt
  $ echo "Test file 2 contents" > mydata/file2.txt
  $ mkbundle -v -o=mybundle.go mydata
  mkbundle: Generating mybundle.go
  mkbundle: + file1.txt
  mkbundle: + file2.txt
```

As a result of running the command above, the file "mybundle.go" will
be generated and it will contain the following Go code:

```Go
  // Bundle file
  // Auto-generated. !! DO NOT EDIT !!
  // Generated: 2013-12-29T22:46:58+02:00

  package main

  import "github.com/npat-efault/bundle"

  var _bundle = []bundle.Entry{
  { Name : "file1.txt",
    Size : 21,
    Gzip : false,
    Data : `
  VGVzdCBmaWxlIDEgY29udGVudHMK
  `},
  { Name : "file2.txt",
    Size : 21,
    Gzip : false,
    Data : `
  VGVzdCBmaWxlIDIgY29udGVudHMK
  `},
  }

  var _bundleIdx bundle.Index

  func init() {
        _bundleIdx = bundle.MkIndex(_bundle)
  }

  // End of bundle
```

This generated file can be linked to your programm allowing access to
the embedded data. The bundle module (github.com/npat-efault/bundle)
contains functions that help you access the data that have been
embedded in the binary by "mkbundle".

## Install

Say:

```
  $ go get github.com/npat-efault/bundle
  $ cd $GOPATH/src/github.com/npat-efault/bundle
  $ ./all.sh install
```

Test with

```
  $ ./all.sh test [-v]
```

Dircetory "serveb" contains an example program. It implements a simple
server that serves bundle data over HTTP. To build and run say:

```
  $ cd serveb
  $ ./all.sh build
  $ ./serveb 6060
```

Then direct your browser to http://localhost:6060/
