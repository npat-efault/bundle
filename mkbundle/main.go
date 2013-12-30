// mkbundle main

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func emitBundleHeader(w io.Writer, pkg, bundle, index string) error {
	var err error
	_, err = fmt.Fprintf(w, BundleHeadFormat,
		pkg, bundle, index,
		BundleImportPath,
		time.Now().Format(time.RFC3339))
	return err
}

func emitBundleFooter(w io.Writer, bundle, index string) error {
	var err error
	_, err = fmt.Fprintf(w, BundleFootFormat, bundle, index)
	return err
}

func emitFile(w io.Writer, fpath, name string, sz int, zip bool) error {
	var f *os.File
	var gw io.WriteCloser
	var err error

	f, err = os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if zip {
		gw, err = NewGoZipWriter(w, name, sz)
	} else {
		gw, err = NewGoWriter(w, name, sz)
	}
	if err != nil {
		return err
	}
	defer gw.Close()
	_, err = io.Copy(gw, f)

	return err
}

func walkDir(w io.Writer, fpath string) error {
	// Walk directory
	var wf = func(p string, i os.FileInfo, e error) error {
		var nm string
		var err error
		var ok bool

		if e != nil {
			return e
		}
		// Handle skip-patterns
		for _, pat := range fl.skip {
			ok, _ = filepath.Match(pat, i.Name())
			if ok {
				if i.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		if i.IsDir() {
			// continue
			return nil
		} else if i.Mode().IsRegular() {
			nm, err = filepath.Rel(fpath, p)
			if err != nil {
				return err
			}
			if fl.verbose {
				log.Printf("+ %s", nm)
			}
			return emitFile(w, p, nm,
				int(i.Size()), fl.gzip)
		} else {
			log.Printf("%s: skipped non-regular file", p)
			return nil
		}
	}
	return filepath.Walk(fpath, wf)
}

func emitBundle(w io.Writer, fpath string) error {
	var info os.FileInfo
	var err error

	err = emitBundleHeader(w, fl.pkg, fl.bundle, fl.index)
	if err != nil {
		return err
	}

	info, err = os.Lstat(fpath)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		// Emit signle file
		name := path.Base(fpath)
		err = emitFile(w, fpath, name, int(info.Size()),
			fl.gzip)
		if err != nil {
			return err
		}
	} else if info.Mode().IsDir() {
		// Walk subtree rooted at dir
		err = walkDir(w, fpath)
		if err != nil {
			return err
		}
	} else {
		// Oops!
		err = fmt.Errorf("%s: not a regular file or directory",
			fpath)
		return err
	}

	err = emitBundleFooter(w, fl.bundle, fl.index)
	if err != nil {
		return err
	}
	return nil
}

func isYounger(ofn string, ifn string) bool {
	var oinf os.FileInfo
	var err error

	if ofn == "" {
		return false
	}
	oinf, err = os.Stat(ofn)
	if err != nil || !oinf.Mode().IsRegular() {
		return false
	}

	var wf = func(p string, iinf os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !oinf.ModTime().After(iinf.ModTime()) {
			err := errors.New("this is younger")
			return err
		}
		return nil
	}
	err = filepath.Walk(ifn, wf)
	if err != nil {
		return false
	}
	return true
}

func main() {
	var fo *os.File
	var err error

	flag.Parse()
	if fl.help {
		flag.CommandLine.SetOutput(os.Stdout)
		fmt.Printf(usage, path.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Println()
		return
	}
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr,
			"incorrect number of arguments.\n")
		flag.Usage()
		os.Exit(1)
	}
	if !fl.always && isYounger(fl.out, flag.Arg(0)) {
		if fl.verbose {
			log.Printf("%s is younger than %s",
				fl.out, flag.Arg(0))
		}
		return
	}
	if fl.out != "" {
		fo, err = os.Create(fl.out)
		if err != nil {
			log.Fatal(err)
		}
		defer fo.Close()
		if fl.verbose {
			log.Printf("Generating %s", fl.out)
		}
	} else {
		fo = os.Stdout
		if fl.verbose {
			log.Print("Generating on <stdout>")
		}
	}
	err = emitBundle(fo, flag.Arg(0))
	if err != nil {
		if fl.out != "" {
			os.Remove(fl.out)
		}
		log.Fatal(err)
	}
}

// Setup for command line arguments parsing

type patlist []string

func (pl *patlist) String() string {
	return fmt.Sprint(*pl)
}

func (pl *patlist) Set(value string) error {
	*pl = append(*pl, value)
	return nil
}

var fl struct {
	out     string
	pkg     string
	bundle  string
	index   string
	gzip    bool
	skip    patlist
	always  bool
	verbose bool
	help    bool
}

func init() {
	flag.Var(&fl.skip, "skip", "Files/dirs to skip (glob pattern)")
	flag.StringVar(&fl.out, "out", "",
		"Output file (if empty, use <stdout>)")
	flag.StringVar(&fl.out, "o", "",
		"Short for \"-out\"")
	flag.StringVar(&fl.pkg, "pkg", "main",
		"Package for the generated source file")
	flag.StringVar(&fl.bundle, "bundle", "_bundle",
		"Name of global that keeps embedded data")
	flag.StringVar(&fl.index, "index", "_bundleIdx",
		"Name of global filename-to-data index")
	flag.BoolVar(&fl.gzip, "g", false,
		"Short for '-gzip'")
	flag.BoolVar(&fl.gzip, "gzip", false,
		"Compress data before embedding")
	flag.BoolVar(&fl.always, "always", false,
		"Regenerate output even if younger than input")
	flag.BoolVar(&fl.always, "a", false,
		"Short for \"-always\"")
	flag.BoolVar(&fl.verbose, "verbose", false,
		"Print actions performed on <stderr>")
	flag.BoolVar(&fl.verbose, "v", false,
		"Short for \"-verbose\"")
	flag.BoolVar(&fl.help, "help", false,
		"Show instructions")
	flag.BoolVar(&fl.help, "h", false,
		"Short for \"-help\"")
	flag.Usage = func() {
		log.Printf("run with '-help' for instructions")
	}
	log.SetFlags(0)
	log.SetPrefix(path.Base(os.Args[0]) + ": ")
	log.SetOutput(os.Stderr)
}
