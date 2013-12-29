package bundle

type Entry struct {
	Name string
	Size int
	Gzip bool
	Data string
}

type Index map[string]*Entry

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

func (idx Index) Has(name string) bool {
	return false
}

func (idx Index) Entry(name string) (*Entry, bool) {
	return nil, false
}

func (idx Index) Size(name string) (int, bool) {
	return 0, false
}

func (idx Index) Gzip(name string) (bool, bool) {
	return false, false
}

func (idx Index) Data(name string) (string, bool) {
	return nil, false
}
