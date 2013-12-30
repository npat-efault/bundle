/*

Command mkbundle embeds arbitrary data files in a Go binary.

Command "mkbundle" allows, moderately sized, arbitrary data files to be
embedded (bundled) inside a Go binary. It generates a Go source file
containing variable initializations with data from the files you wish
to embed. Data are embedded in base64 encoding. This generated file
can be compiled and linked with the rest of your program's code. Usage is:

  mkbundle [flags] <file-or-dir>

The following flags are recognized:

  -a=false: Short for "-always"
  -always=false: Regenerate output even if younger than input
  -bundle="_bundle": Name of global that keeps embedded data
  -g=false: Short for '-gzip'
  -gzip=false: Compress data before embedding
  -h=false: Short for "-help"
  -help=false: Show instructions
  -index="_bundleIdx": Name of global filename-to-data index
  -o="": Short for "-out"
  -out="": Output file (if empty, use <stdout>)
  -pkg="main": Package for the generated source file
  -skip=[]: Files/dirs to skip (glob pattern)
  -v=false: Short for "-verbose"
  -verbose=false: Print actions performed on <stderr>

The <file-or-dir> argument is the name (path) of the file you wish to
embed. If a directory name is given instead, all files in that
directory (and its subdirectories, recursively) will be embedded. The
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
used to map file-names to embedded data. File names in the global
index are made relative to the direcory given as argument to the
mkbundle command. If a single file (not directory) name is given as
argument to the command, then the file-name in the index will be the
base-name of that single file.

If the '-gzip' flag is given, then files will be compressed with gzip
before being embedded.

If the '-verbose' flag is given, then the command will print a few
short messages on <stderr> indicating the actions it performs. Without
'-verbose' the commands prints messages only on errors, otherwise it
remains completely silent.

If the output file (specified by the "-out" flag) already exists, it
will be re-generated only if <file-or-dir>, or at least one of the
files and sub-directories in it, are younger than the output file. You
can override this behavior using the "-always" flag.

For information on how to access the bundled data from your code, see
the documentation of package:

  github.com/npat-efault/bundle

*/
package main
