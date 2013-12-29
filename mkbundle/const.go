// Constants and strings used by mkbundle

package main

const usage = ` 
Usage is: %[1]s [flags] <file-or-dir> 

Command "%[1]s" allows, moderately sized, abitrary data files to be
embedded (bundled) inside a Go binary. It generates a Go source file
containing variable initializations with data from the files you wish
to embed. Data are embedded in base64 encoding. This generated file
can be compiled and linked with the rest of your program's code.

The <file-or-dir> argument is the name (path) of the file you wish to
embed. If a directory name is given instead, all files in that
directory (and its subridectories, recursivelly) will be embedded. The
'-skip' flag can be used to specify files and directories that will be
skipped when generating the bundle. The argument of the '-skip' flag
is interpreted as a glob pattern. The '-skip' flag can be given
multiple times, if files and directories matching multiple patterns
must be skipped.

The '-pkg' flag provides the name of the package the generated file
will belong to. The '-bundle' flag provides the name of the global
variable that will be defined in the generated file to reference the
embedded data. The '-index' flag provides the name of the global map
variable that will be defined in the generated file. This variable is
used to map file-names to embedded data.

If the output file (specified by the "-out" flag) already exists, it
will be re-generated only if <file-or-dir>, or at least one of the
files and sub-directories in it, are younger than the output file. You
can override this behavior using the "-always" flag.

For information on how to access the bundled data from your code, see
the documentation of package:

  %[2]s

The following flags are recognized by %[1]s:

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
