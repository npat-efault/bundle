package bundle_test

import (
	"bytes"
	"compress/gzip"
	"github.com/npat-efault/bundle"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var data_dir = "test_data/"

func mkentries(data_dir string) ([]string, error) {
	var entries []string
	var err error

	var wf = func(p string, inf os.FileInfo, e error) error {
		var err error
		var nm string
		if e != nil {
			return e
		}
		if !inf.Mode().IsRegular() {
			return nil
		}
		nm, err = filepath.Rel(data_dir, p)
		if err != nil {
			return err
		}
		entries = append(entries, nm)
		return nil
	}
	err = filepath.Walk(data_dir, wf)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func TestIndex(t *testing.T) {
	var entries []string
	var d bundle.Dir
	var e *bundle.Entry
	var nm string
	var i, sz int
	var err error

	entries, err = mkentries(data_dir)
	if err != nil {
		t.Fatalf("mkentries failed (%s): %s", data_dir, err)
	}
	if sz = len(_bundleIdx); sz != len(entries) {
		t.Fatalf("Index size %d != %d", sz, len(entries))
	}
	if sz = len(_bundleIdx.Dir("")); sz != len(entries) {
		t.Fatalf("Dir returned %d != %d", sz, len(entries))
	}
	for _, nm = range entries {
		if !_bundleIdx.Has(nm) {
			t.Fatalf("Not in index: %s", nm)
		}
	}
	for i = 0; i < len(entries); i++ {
		d = _bundleIdx.Dir(entries[i])
		e = _bundleIdx.Entry(entries[i])
		if e == nil {
			t.Fatal("Cannot find entry: %s", entries[0])
		}
		// A name may be a prefix of another name!
		if len(d) < 1 {
			t.Fatalf("Cannot dir entry: %s", entries[0])
		}
		if e.Name != d[0].Name || e.Size != d[0].Size {
			t.Fatalf("Oops! Name: %s (%s), Size %d (%d)",
				e.Name, d[0].Name,
				e.Size, d[0].Size)
		}
		t.Logf("Entry: %s, Size: %d, Gzip %v",
			e.Name, e.Size, e.Gzip)
	}
}

func TestData(t *testing.T) {
	var entries []string
	var e *bundle.Entry
	var br *bundle.Reader
	var data, rdata, fdata []byte
	var i int
	var err error

	entries, err = mkentries(data_dir)
	if err != nil {
		t.Fatalf("mkentries failed (%s): %s", data_dir, err)
	}
	for i = 0; i < len(entries); i++ {
		// Read data from file
		fdata, err = ioutil.ReadFile(data_dir + entries[i])
		if err != nil {
			t.Fatalf("ReadFile(): %s", err)
		}
		// Get bundle entry
		e = _bundleIdx.Entry(entries[i])
		if e == nil {
			t.Fatalf("Entry not found: %s", entries[i])
		}
		// Get data from bundle using Decode
		data, err = bundle.Decode(e, 0)
		if err != nil {
			t.Fatalf("bundle.Decode(): %s", err)
		}
		if len(data) != e.Size {
			t.Fatalf("len(data) %d != e.Size %d",
				len(data), e.Size)
		}
		if len(data) != len(fdata) {
			t.Fatalf("Bad data sz for: %s", entries[i])
		}
		if bytes.Compare(data, fdata) != 0 {
			t.Fatalf("Bad data for: %s", entries[i])
		}
		// Get data from bundle using bundle.Reader
		br, err = _bundleIdx.Open(entries[i], 0)
		if err != nil {
			t.Fatalf("idx.Open(): %s", err)
		}
		rdata, err = ioutil.ReadAll(br)
		if err != nil {
			t.Fatalf("ReadAll(br): %s", err)
		}
		if len(rdata) != e.Size {
			t.Fatalf("len(rdata) %d != e.Size %d",
				len(rdata), e.Size)
		}
		if len(rdata) != len(fdata) {
			t.Fatalf("Bad rdata sz for: %s", entries[i])
		}
		if bytes.Compare(rdata, fdata) != 0 {
			t.Fatalf("Bad rdata for: %s", entries[i])
		}
		t.Logf("Entry: %s, Size: %d, Gzip: %v",
			e.Name, e.Size, e.Gzip)
	}
}

func TestCompressed(t *testing.T) {
	var entries []string
	var e *bundle.Entry
	var br *bundle.Reader
	var gr *gzip.Reader
	var data, rdata, fdata, ddata []byte
	var i int
	var err error

	entries, err = mkentries(data_dir)
	if err != nil {
		t.Fatalf("mkentries failed (%s): %s", data_dir, err)
	}
	for i = 0; i < len(entries); i++ {
		// Read data from file
		fdata, err = ioutil.ReadFile(data_dir + entries[i])
		if err != nil {
			t.Fatalf("ReadFile error: %s", err)
		}
		// Get bundle entry
		e = _bundleIdx.Entry(entries[i])
		if e == nil {
			t.Fatalf("Not found: %s", entries[i])
		}
		// Get data from bundle using Decode
		data, err = bundle.Decode(e, bundle.NODC)
		if err != nil {
			t.Fatalf("bundle.Decode(): %s", err)
		}
		gr, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("gzip.NewReader: %s", err)
		}
		ddata, err = ioutil.ReadAll(gr)
		if err != nil {
			t.Fatalf("ReadAll(gr): %s")
		}
		err = gr.Close()
		if err != nil {
			t.Fatalf("gr.Close(): %s")
		}
		if len(ddata) != e.Size {
			t.Fatalf("len(ddata) %d != e.Size %d",
				len(ddata), e.Size)
		}
		if len(ddata) != len(fdata) {
			t.Fatalf("Bad ddata sz for: %s", entries[i])
		}
		if bytes.Compare(ddata, fdata) != 0 {
			t.Fatalf("Bad ddata for: %s", entries[i])
		}
		// Get data from bundle using bundle.Reader
		br, err = _bundleIdx.Open(entries[i], bundle.NODC)
		if err != nil {
			t.Fatalf("bundle.Open(): %s", err)
		}
		gr, err = gzip.NewReader(br)
		if err != nil {
			t.Fatalf("gzip.NewReader(br): %s")
		}
		ddata, err = ioutil.ReadAll(gr)
		if err != nil {
			t.Fatalf("ReadAll(gr): %s")
		}
		err = gr.Close()
		if err != nil {
			t.Fatalf("gr.Close(): %s")
		}
		if len(ddata) != e.Size {
			t.Fatalf("len(ddata) %d != e.Size %d",
				len(rdata), e.Size)
		}
		if len(ddata) != len(fdata) {
			t.Fatalf("Bad ddata sz for: %s", entries[i])
		}
		if bytes.Compare(ddata, fdata) != 0 {
			t.Fatalf("Bad ddata for: %s", entries[i])
		}
		t.Logf("Entry: %s, Size: %d, Gzip: %v",
			e.Name, e.Size, e.Gzip)
	}
}
