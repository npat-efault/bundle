// Constants and strings used by mkbundle

package main

const usage = ` 
Usage is: %[1]s [flags] <file-or-dir> 

Command "%[1]s" allows, moderately sized, abitrary data files to be
embedded (bundled) inside a Go binary. It generates a Go source file
containing variable initializations with data from the files you wish
to embed in base64 encoding. The generated file can be compiled and
linked with the rest of your program's code.

The <file-or-dir> argument is the name (path) of the file you wish to
embed. If a directory name is given instead, all files in that
directory (and its subridectories, recursivelly) will be embedded.

If the output file (specified by the "-out" flag) already exists, it
will be re-generated only if <file-or-dir>, or at least one of the
files and sub-directories in it, are younger than the output file. You
can override this behavior using the "-always" flag.

The following flags are recognized:

`

const BundleImportPath string = "github.com/npat-efault/bundle"

const BundleHeadFormat string = `
// Bundle file
// Auto-generated. !! DO NOT EDIT !!
// Generated: %[5]s

package %[1]s

import "%[4]s"

var %[2]s = []bundle.Entry{
`

const BundleFootFormat string = `}

var %[2]s bundle.Index

func Init() {
     %[2]s = bundle.MkIndex(%[1]s)
}

// End of bundle
`

const FileHeadFormat string = `{ Name : "%[1]s",
  Size : %[2]d,
  Gzip : %[3]v,
  Data : ` + "`"

const FileFootFormat string = "\n`},\n"
