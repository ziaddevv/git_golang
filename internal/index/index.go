package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"slices"
)

type Index struct {
	Entries []*Entry
}

func (idx *Index) Pack() []byte {
	buf := new(bytes.Buffer)

	// header
	buf.WriteString("DIRC")
	binary.Write(buf, binary.BigEndian, uint32(2))                // version
	binary.Write(buf, binary.BigEndian, uint32(len(idx.Entries))) // count

	// entries
	for _, e := range idx.Entries {
		buf.Write(e.Pack())
	}

	// sha1 checksum of everything so far
	sum := sha1.Sum(buf.Bytes())
	buf.Write(sum[:])

	return buf.Bytes()
}

func (idx *Index) Add(e *Entry) {
	for i := range idx.Entries {
		if e.Path == idx.Entries[i].Path {
			// Update the actual slice element
			idx.Entries[i] = e
			return
		}
	}

	// Not found, so append
	idx.Entries = append(idx.Entries, e)
}

func (idx *Index) Remove(path string) {
	idx.Entries = slices.DeleteFunc(idx.Entries, func(e *Entry) bool {
		return e.Path == path
	})
}

// func (idx *Index) remove(path string) {
//     n := 0
//     for _, entry := range idx.Entries {
//         if entry.Path != path {
//             idx.Entries[n] = entry
//             n++
//         }
//     }
//     idx.Entries = idx.Entries[:n]
// }

func UnpackIndex(data []byte) *Index {
	// skip header (12 bytes)
	count := binary.BigEndian.Uint32(data[8:12])

	idx := &Index{}
	offset := 12

	for i := 0; i < int(count); i++ {
		entry, newOffset := UnpackEntry(data, offset)
		idx.Entries = append(idx.Entries, entry)
		offset = newOffset
	}

	return idx
}
