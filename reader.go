package diskbst

import (
	"bytes"
	"errors"
	"os"
	"syscall"
)

type reader struct {
	data []byte
}

type Reader interface {
	Get(key []byte) ([]byte, error)
	Close()
}

func OpenReader(pathName string) (Reader, error) {
	var r reader

	f, err := os.Open(pathName)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// validate magic number
	mn := make([]byte, len(magicNumber))
	_, err = f.Read(mn)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(mn, magicNumber) != 0 {
		return nil, errInvalidMagicNumber
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	d, err := syscall.Mmap(int(f.Fd()), 0, int(info.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}

	r.data = d

	return &r, nil
}

func (r *reader) Get(key []byte) ([]byte, error) {
	if len(r.data) == len(magicNumber) {
		return nil, errKeyNotFound
	}

	var currNode node

	// skip node size as it is not utilized by the reader
	currPos := uint64(len(magicNumber))
	currPos += 8

	for currPos < uint64(len(r.data)) {
		currNode.deserialize(r.data[currPos:])

		if bytes.Compare(key, currNode.key) == 0 {
			return currNode.value, nil
		}

		var next uint64

		if bytes.Compare(key, currNode.key) <= 0 {
			next = currNode.leftChild
		}

		if bytes.Compare(key, currNode.key) > 0 {
			next = currNode.rightChild
		}

		if next == 0 {
			return nil, errKeyNotFound
		}

		currPos = next

		// skip node size as it is not utilized by the reader
		currPos += 8
	}

	return nil, errKeyNotFound
}

func (r *reader) Close() {
	_ = syscall.Munmap(r.data)
}

var errKeyNotFound = errors.New("key not found")
