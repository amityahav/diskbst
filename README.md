## Disk based binary search tree
just a simple disk serialized binary search tree.

It can be used in an LSM kv store for flushing the memtable to disk as a binary search tree.

### Example usage:

``` go
package main

import (
	"fmt"
	diskbst "github.com/amityahav/disk_bst"
)

func main() {
	memtable := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	// Open a new file to store the tree
	writer, _ := diskbst.OpenWriter("/path/to/bst")

	// Store all key-value pairs in the tree
	for key, val := range memtable {
		_ = writer.Put([]byte(key), val)
	}

	writer.Close()

	// Open a reader for the tree
	// the tree file will be memory-mapped
	reader, _ := diskbst.OpenReader("/path/to/bst")

	val, _ := reader.Get([]byte("key1"))
	fmt.Printf(string(val)) // value1

	_, err := reader.Get([]byte("doesnt_exist"))
	fmt.Printf(err.Error()) // key not found
}

```