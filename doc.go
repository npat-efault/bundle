/*

Package bundle: Access data that have been embedded in a Go binary by
mkbundle.

Package bundle, together with the "mkbundle" command, allow, moderately
sized, arbitrary data files to be embedded (bundled) inside a Go
binary.

The mkbundle command generates a Go source file that contains global
variables initialized with data from the files you wish to embed,
(optionally) compressed with gzip, and encoded with base64. Example:

  $ mkdir mydata
  $ echo "Test file 1 contents" > mydata/file1.txt
  $ echo "Test file 2 contents" > mydata/file2.txt
  $ mkbundle -v -o=mybundle.go mydata
  mkbundle: Generating mybundle.go
  mkbundle: + file1.txt
  mkbundle: + file2.txt

As a result of running the command above, the file "mybundle.go" will
be generated and it will contain the following Go code:

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

As you can see, a global variable named "_bundle" is defined which is
a slice with one entry for each of the files included in the
bundle. Each entry keeps the file's name, it's size (the original
size, before compression and encoding), an indication whether the file
was compressed, and the file's data in base64 encoding. In addition a
global map, named "_bundleIdx" is defined which associates file-names
with the bundle entries. This generated file can be linked to your
programm allowing access to the embedded data. Assume this code in a
file "tstbundle.go" :

  package main
  import "fmt"

  func main() {
      for _, e := range _bundleIdx.Dir("") {
          b := bundle.Decode(e, 0)
          fmt.Println(e.Name, e.Size, string(b))
      }
  }

The Dir() method, defined on _bundleIdx, returns a slice of pointers
to the entries in the index with names matching the given prefix (all
entries for an empty prefix) sorted by name in ascending order.

The code above, compiled and linked together with the generated file
"mybundle.go", when run produces the output:

  $ tstbundle
  file1.txt 21 Test file 1 contents
  file2.txt 21 Test file 2 contents

The name of the global variable keeping the embedded data ("_bundle"
in our example), the name of the global name-to-entry index
("_bundleIdx"), the name of the package for the generated bundle
("main"), whether to compress the files or not, and several other
options can be controlled by flags passed to the "mkbundle"
command. See the command's documentation for more information.

Summarizing: The command "mkbundle" allows arbitrary data files to be
embedded in Go binaries by converting the files to statements
initializing global variables. This module
("github.com/npat-efault/bundle") provides an interface that can be
used to access the data embedded in the binary.

*/
package bundle
