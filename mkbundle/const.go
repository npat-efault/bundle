// Constants and strings used by mkbundle

package main

const usage = ` 
Usage is: %[1]s [flags] <file-or-dir> 

Command "%[1]s" allows, moderately sized, arbitrary data files to be
embedded (bundled) inside a Go binary. 

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

func init() {
     %[2]s = bundle.MkIndex(%[1]s)
}

// End of bundle
`

const FileHeadFormat string = `{ Name : "%[1]s",
  Size : %[2]d,
  Gzip : %[3]v,
  Data : ` + "`"

const FileFootFormat string = "\n`},\n"
