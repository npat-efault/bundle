// Interface for accessing the bundle

package bundle

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Falgs for Index.Open and Decode
const (
	NODC int = 1 << iota // Do not decompress data
)

// The Entry type is used to represent the bundled file data. One such
// structure is used for every bundled file. All structures are kept
// in a global slice. When files are included in a bundle the terms
// "entry" and "file" may be used interchangeably (i.e. the bundle
// contains 5 files / entries).
type Entry struct {
	// The name of the entry
	Name string
	// The size of the entry in bytes (size of "Data"). This the
	// original data size, before compression and encoding.
	Size int
	// Is the entry compressed?
	Gzip bool
	// Entry data compressed (if Gzip is true) and base64 encoded
	Data string
}

// Index is the type of the global map of names to entries. Such a map
// is declared for every bundle.
type Index map[string]*Entry

// MkIndex creates and initializes the names-to-entries index. A call
// to MkIndex is inserted automatically by "mkbundle" to the "init"
// function of the generated file.
func MkIndex(bundle []Entry) Index {
	var bsz int
	var idx Index

	bsz = len(bundle)
	idx = make(Index, bsz)
	for i := 0; i < bsz; i++ {
		idx[bundle[i].Name] = &bundle[i]
	}
	return idx
}

// The Has method returns true if the bundle has an entry with the
// given name
func (idx Index) Has(name string) bool {
	_, ok := idx[name]
	return ok
}

// The Entry method returns a pointer to the entry with the requested
// name, if such an entry exists in the bundle, or nil if no such
// entry exists.
func (idx Index) Entry(name string) *Entry {
	e, ok := idx[name]
	if !ok {
		return nil
	}
	return e
}

// Dir is a slice of pointers to entries. It implements sort.Interface
type Dir []*Entry

// Len returns length of Dir (number of elements)
func (d Dir) Len() int {
	return len(d)
}

// Less returns d[i].Name < d[j].Name
func (d Dir) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

// Swap swaps entries i and j in Dir
func (d Dir) Swap(i, j int) {
	t := d[i]
	d[i] = d[j]
	d[j] = t
}

// The Dir method returns a Dir (slice of Entry pointers) of all the
// entries whose names match the given prefix (all entries whose names
// start with string "prefix")
func (idx Index) Dir(prefix string) []*Entry {
	var dir Dir
	for _, e := range idx {
		if strings.HasPrefix(e.Name, prefix) {
			dir = append(dir, e)
		}
	}
	sort.Sort(dir)
	return dir
}

// TODO(npat): Change Decode() to allow return of compressed data
// (optionally)

// Decode returns the decoded data for the bundle entry pointed to by
// "e". Returns a slice of bytes with the decoded, decompressed (if
// required), ready to use entry data, and an error indication which
// is not-nil if the data cannot be decoded. If argument "flag" is
// NODC, and the entry data are compressed (Entry.Gzip == true),
// Decode will not decompress the data it returns (it will only decode
// them). In most cases it is preferable to use the Reader interface
// instead of calling Decode.
func Decode(e *Entry, flag int) ([]byte, error) {
	var rs *strings.Reader
	var r64 io.Reader
	var rz *gzip.Reader
	var buf *bytes.Buffer
	var err error

	rs = strings.NewReader(e.Data)
	r64 = base64.NewDecoder(base64.StdEncoding, rs)
	if e.Gzip && (flag&NODC == 0) {
		rz, err = gzip.NewReader(r64)
		if err != nil {
			return nil, err
		}
		defer rz.Close()
	} else {
		rz = nil
	}
	buf = new(bytes.Buffer)
	if rz != nil {
		_, err = io.Copy(buf, rz)
	} else {
		_, err = io.Copy(buf, r64)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TODO(npat): Change NewReader to allow return of compressed data
// (optionally)

// A Reader implents the io.Reader and io.Closer interface by reading,
// decoding, and decompressing (if required) data from a bundle entry.
type Reader struct {
	rs  *strings.Reader
	r64 io.Reader
	rz  *gzip.Reader
}

// Open intializes and returns a Reader that reads from the bundle
// entry specified by "name". It returns an error if no such entry
// exists, or if the reader cannot be initialized. If argument "flag"
// is NODC, and the entry data are compressed (Entry.Gzip == true),
// the reader will not decompress the data read from it (it will only
// decode them).
func (idx Index) Open(name string, flag int) (*Reader, error) {
	var entry *Entry
	var br *Reader
	var err error

	if !idx.Has(name) {
		err = fmt.Errorf("no such entry: %s", name)
		return nil, err
	}
	entry = idx.Entry(name)
	br = &Reader{}
	br.rs = strings.NewReader(entry.Data)
	br.r64 = base64.NewDecoder(base64.StdEncoding, br.rs)
	if entry.Gzip && (flag&NODC == 0) {
		br.rz, err = gzip.NewReader(br.r64)
		if err != nil {
			return nil, err
		}
	} else {
		br.rz = nil
	}
	return br, nil
}

// The Read method is used to read data from a bundle entry. Read
// fills slice "p" with decoded, decompressed, ready to use
// data. Returns the number of bytes read (stored in "p") and an error
// indication (which is not-nil when a read error has occured).
func (br *Reader) Read(p []byte) (int, error) {
	if br.rz != nil {
		return br.rz.Read(p)
	} else {
		return br.r64.Read(p)
	}
}

// The Close method is Used to terminate the operation of the
// Reader. Returns an error indication which is not-nil if an error
// has occured during close. After calling Close no other operations
// must be performed on this Reader.
func (br *Reader) Close() error {
	if br.rz != nil {
		return br.rz.Close()
	} else {
		return nil
	}
}
