package diskbst

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"
)

func TestBST(t *testing.T) {
	t.Run("write/read", func(t *testing.T) {
		path := t.TempDir()
		pathName := filepath.Join(path, "test.bst")

		w, err := OpenWriter(pathName)
		if err != nil {
			t.Error(err)
			return
		}

		err = w.Put([]byte("hello"), []byte("world"))
		if err != nil {
			t.Error(err)
			return
		}

		err = w.Put([]byte("helloo"), []byte("worldd"))
		if err != nil {
			t.Error(err)
			return
		}

		err = w.Put([]byte("hell"), []byte("worl"))
		if err != nil {
			t.Error(err)
			return
		}

		err = w.Put([]byte("hellooo"), []byte("worlddd"))
		if err != nil {
			t.Error(err)
			return
		}

		w.Close()

		r, err := OpenReader(pathName)
		if err != nil {
			t.Error(err)
			return
		}

		val, err := r.Get([]byte("hello"))
		if err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(val, []byte("world")) != 0 {
			t.Fail()
			return
		}

		val, err = r.Get([]byte("hellooo"))
		if err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(val, []byte("worlddd")) != 0 {
			t.Fail()
			return
		}

		val, err = r.Get([]byte("hell"))
		if err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(val, []byte("worl")) != 0 {
			t.Fail()
			return
		}

		val, err = r.Get([]byte("blah"))
		if !errors.Is(err, errKeyNotFound) {
			t.Fail()
			return
		}

	})
}
