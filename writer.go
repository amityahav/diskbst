package diskbst

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type writer struct {
	cursor *os.File
	tail   int64
}

type Writer interface {
	Put(key []byte, value []byte) error
	Close()
}

func OpenWriter(pathName string) (Writer, error) {
	var w writer

	cursor, err := os.OpenFile(pathName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	w.cursor = cursor

	info, err := cursor.Stat()
	if err != nil {
		return nil, err
	}

	w.tail = info.Size()

	return &w, nil
}

func (w *writer) Put(key []byte, value []byte) error {
	newNode := node{
		key:   key,
		value: value,
	}

	pos, err := w.findPos(key)
	if err != nil {
		return err
	}

	// write new node to disk
	data := newNode.serialize()
	n, err := w.cursor.WriteAt(data, w.tail)
	if err != nil {
		return err
	}

	if n < len(data) {
		return fmt.Errorf("failed writing new node")
	}

	newTail := w.tail + int64(n)

	if isRoot := pos == 0; isRoot {
		w.tail = newTail
		return nil
	}

	// link new node to its parent
	nt := make([]byte, 8)
	binary.LittleEndian.PutUint64(nt, uint64(w.tail))
	n, err = w.cursor.WriteAt(nt, pos)
	if err != nil {
		return err
	}

	if n < len(nt) {
		return fmt.Errorf("failed writing new node offset to parent")
	}

	w.tail = newTail

	return nil
}

func (w *writer) Close() {
	_ = w.cursor.Close()
}

func (w *writer) findPos(key []byte) (int64, error) {
	var (
		currPos        int64
		currNode       node
		childPtrOffset int64
	)
	nodeSize := make([]byte, 8)

	if w.tail == 0 {
		return currPos, nil
	}

	for {
		n, err := w.cursor.ReadAt(nodeSize, currPos)
		if err != nil {
			return 0, err
		}

		if n < len(nodeSize) {
			// corruption
			return 0, fmt.Errorf("failed reading node size")
		}

		currPos += 8

		s := binary.LittleEndian.Uint64(nodeSize)

		nd := make([]byte, s)
		n, err = w.cursor.ReadAt(nd, currPos)
		if err != nil {
			return 0, err
		}

		if n < len(nd) {
			// corruption
			return 0, fmt.Errorf("failed reading node")
		}

		var next int64
		currNode.deserialize(nd)

		if bytes.Compare(key, currNode.key) <= 0 {
			// key <= currNode.Key
			next = int64(currNode.leftChild)
			childPtrOffset = currPos + int64(8+len(currNode.key)+8+len(currNode.value))
		} else {
			// key > currNode.key
			next = int64(currNode.rightChild)
			childPtrOffset = currPos + int64(8+len(currNode.key)+8+len(currNode.value)+8)
		}

		if next == 0 {
			// parent of the new node found
			break
		}

		currPos = next
	}

	return childPtrOffset, nil
}
