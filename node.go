package diskbst

import "encoding/binary"

type node struct {
	key        []byte
	value      []byte
	leftChild  uint64
	rightChild uint64
}

func (n *node) serialize() []byte {
	res := make([]byte, 8+uint64(8+len(n.key)+8+len(n.value)+8+8))

	var pos int

	// node length
	binary.LittleEndian.PutUint64(res[pos:], uint64(8+len(n.key)+8+len(n.value)+8+8))
	pos += 8

	// key length
	binary.LittleEndian.PutUint64(res[pos:], uint64(len(n.key)))
	pos += 8

	// key
	copy(res[pos:pos+len(n.key)], n.key)
	pos += len(n.key)

	// value length
	binary.LittleEndian.PutUint64(res[pos:], uint64(len(n.value)))
	pos += 8

	// value
	copy(res[pos:pos+len(n.value)], n.value)
	pos += len(n.value)

	// left child
	binary.LittleEndian.PutUint64(res[pos:], n.leftChild)
	pos += 8

	// right child
	binary.LittleEndian.PutUint64(res[pos:], n.rightChild)
	pos += 8

	return res
}

func (n *node) deserialize(data []byte) {
	var pos int

	keyLen := binary.LittleEndian.Uint64(data[pos:])
	pos += 8
	n.key = make([]byte, keyLen)
	copy(n.key, data[pos:])
	pos += len(n.key)

	valueLen := binary.LittleEndian.Uint64(data[pos:])
	n.value = make([]byte, valueLen)
	pos += 8
	copy(n.value, data[pos:])
	pos += len(n.value)

	n.leftChild = binary.LittleEndian.Uint64(data[pos:])
	pos += 8

	n.rightChild = binary.LittleEndian.Uint64(data[pos:])
	pos += 8
}
