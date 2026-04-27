package index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

// fixed 62-byte portion of each entry
type EntryFixed struct {
	CtimeSec  uint32   // 4 --> when the files's metadata last changes
	CtimeNano uint32   // 4  --> nanosecond part of ctime
	MtimeSec  uint32   // 4  --> when the filecontent last changed
	MtimeNano uint32   // 4	---> nanosecond part of mtime
	Dev       uint32   // 4  ---> which device (disk/partition)the file  lives
	Ino       uint32   // 4   ---> inode number - the files's unique id on that device --> used to detect that file was replaced
	Mode      uint32   // 4  ---> (file type + unix  permissions)
	Uid       uint32   // 4   ---> user id of teh file owner
	Gid       uint32   // 4	  ---> group id of the file owner
	Size      uint32   // 4   -->file size i nbytes on disk
	SHA       [20]byte // 20  --> binary sha1 of the blob object     --> referes to the actual blob object in /objects
	Flags     uint16   // 2  -->
	/*
		bit 15      → assume-unchanged  (git update-index --assume-unchanged)
		bit 14      → extended          (if 1, two more flag bytes follow — version 3 only)
		bits 12-13  → stage             (0=normal, 1/2/3=merge conflict stages)

		i just will store this
		bits 0-11   → name length       (capped at 0xFFF if path > 4095 chars)


	*/
}

// Full entry including variable length path
type Entry struct {
	EntryFixed
	Path string // -->  relative path from repo root  "src/main.go"
}

func (e *Entry) Pack() []byte {
	buf := new(bytes.Buffer)

	// write fixed fields --> big endian
	//binary helps us to convert numbers to bytes
	//binary just writes to bytes.buffer and reads from bytes.reader
	binary.Write(buf, binary.BigEndian, e.EntryFixed)

	// write path + null terminator
	buf.WriteString(e.Path)
	buf.WriteByte(0x00)

	// pad to 8-byte boundary
	// 62 (fixed) + len(path) + 1 (null terminator)
	entryLen := 62 + len(e.Path) + 1
	pad := (8 - (entryLen % 8)) % 8
	for i := 0; i < pad; i++ {
		buf.WriteByte(0x00)
	}
	// fmt.Println(buf.String())

	return buf.Bytes()
}

func UnpackEntry(data []byte, offset int) (*Entry, int) {
	e := &Entry{}

	// read fixed 62 bytes
	fixed := data[offset : offset+62]
	r := bytes.NewReader(fixed)
	binary.Read(r, binary.BigEndian, &e.EntryFixed)
	offset += 62

	// read path until null byte
	end := bytes.IndexByte(data[offset:], 0x00)
	e.Path = string(data[offset : offset+end])
	offset += end + 1

	// skip padding
	entryLen := 62 + len(e.Path) + 1
	pad := (8 - (entryLen % 8)) % 8
	offset += pad

	return e, offset
}

func NewEntryFromFile(path string, sha1Hex string) (*Entry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var sha [20]byte
	raw, err := hex.DecodeString(sha1Hex)
	if err != nil {
		return nil, err
	}
	if len(raw) != 20 {
		return nil, fmt.Errorf("invalid sha1 length")
	}
	copy(sha[:], raw)

	nameLen := len(path)
	if nameLen > 0x0FFF {
		nameLen = 0x0FFF
	}

	var fixed EntryFixed
	fixed.SHA = sha
	fixed.Flags = uint16(nameLen)

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		fixed.CtimeSec = uint32(stat.Ctim.Sec)
		fixed.CtimeNano = uint32(stat.Ctim.Nsec)

		fixed.MtimeSec = uint32(stat.Mtim.Sec)
		fixed.MtimeNano = uint32(stat.Mtim.Nsec)

		fixed.Dev = uint32(stat.Dev)
		fixed.Ino = uint32(stat.Ino)
		// fixed.Mode = uint32(stat.Mode)
		osMode := stat.Mode
		switch {
		case osMode&syscall.S_IFMT == syscall.S_IFLNK:
			fixed.Mode = 0120000
		case osMode&0111 != 0:
			fixed.Mode = 0100755
		default:
			fixed.Mode = 0100644
		}
		fixed.Uid = uint32(stat.Uid)
		fixed.Gid = uint32(stat.Gid)
		fixed.Size = uint32(stat.Size)
	} else {
		sec := uint32(info.ModTime().Unix())

		fixed.CtimeSec = sec
		fixed.CtimeNano = 0
		fixed.MtimeSec = sec
		fixed.MtimeNano = 0
		fixed.Mode = uint32(info.Mode())
		fixed.Size = uint32(info.Size())
	}

	entry := &Entry{
		EntryFixed: fixed,
		Path:       path,
	}

	return entry, nil
}

func NewEntryFromCacheInfo(modeStr, shaHex, path string) (*Entry, error) {

	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid mode: %s", modeStr)
	}

	switch uint32(mode) {
	case 0100644, 0100755, 0120000:
	default:
		return nil, fmt.Errorf("invalid mode: %s", modeStr)
	}

	raw, err := hex.DecodeString(shaHex)
	if err != nil {
		return nil, err
	}
	if len(raw) != 20 {
		return nil, fmt.Errorf("invalid sha1 length")
	}
	var sha [20]byte
	copy(sha[:], raw)

	nameLen := len(path)
	if nameLen > 0x0FFF {
		nameLen = 0x0FFF
	}

	return &Entry{
		EntryFixed: EntryFixed{
			Mode:  uint32(mode),
			SHA:   sha,
			Flags: uint16(nameLen),
			// all stat fields stay zero no file on disk
		},
		Path: path,
	}, nil
}
