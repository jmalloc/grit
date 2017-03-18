package index

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/boltdb/bolt"
)

// Indexer is a function that returns the set of keys applicable to a directory.
type Indexer func(dir string) (keys []string, err error)

// Index is an index of repository locations.
type Index struct {
	db *bolt.DB

	wg  sync.WaitGroup
	err atomic.Value
}

// Open opens the index database at path f.
func Open(f string) (*Index, error) {
	if err := os.MkdirAll(path.Dir(f), 0755); err != nil {
		return nil, err
	}

	db, err := bolt.Open(f, 0644, nil)
	if err != nil {
		return nil, err
	}

	return &Index{db: db}, nil
}

// Close closes the index.
func (i *Index) Close() {
	_ = i.db.Close()
}

// Find returns a list of paths matching the given key.
func (i *Index) Find(key string) (dirs []string, err error) {
	err = i.db.View(func(tx *bolt.Tx) error {
		bucket := readActiveBucket(tx)
		if bucket == nil {
			return nil
		}

		sub := bucket.Bucket([]byte(key))
		if sub == nil {
			return nil
		}

		return sub.ForEach(func(dir []byte, _ []byte) error {
			dirs = append(dirs, string(dir))
			return nil
		})
	})

	return
}

// List returns a slice of all keys that begin with prefix, which may be empty.
func (i *Index) List(prefix string) (keys []string, err error) {
	err = i.db.View(func(tx *bolt.Tx) error {
		bucket := readActiveBucket(tx)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(key []byte, _ []byte) error {
			str := string(key)
			if strings.HasPrefix(str, prefix) {
				keys = append(keys, str)
			}
			return nil
		})
	})

	return
}

// Add indexes a single directory without recursing.
func (i *Index) Add(dir string, fn Indexer) error {
	keys, err := fn(dir)
	if err != nil {
		return err
	}

	return i.db.Update(func(tx *bolt.Tx) error {
		bucket, err := writeActiveBucket(tx)
		if err != nil {
			return err
		}
		return writeKeys(bucket, keys, dir)
	})
}

// Rebuild indexes a set of directol ries recursively, replacing the existing index.
func (i *Index) Rebuild(dirs []string, fn Indexer) error {
	var nextBucket []byte

	// allocate a new bucket that we build the next index into ...
	if err := i.db.Update(func(tx *bolt.Tx) error {
		var err error
		nextBucket, _, err = writeNextBucket(tx)
		return err
	}); err != nil {
		return err
	}

	defer i.gc()

	// walk the directories ...
	for _, dir := range dirs {
		i.wg.Add(1)
		go i.walk(nextBucket, dir, fn)
	}

	i.wg.Wait()
	if err, _ := i.err.Load().(error); err != nil {
		return err
	}

	// set the new bucket as the active one ...
	return i.db.Update(func(tx *bolt.Tx) error {
		return setActiveBucket(tx, nextBucket)
	})
}

func (i *Index) walk(bucket []byte, dir string, fn Indexer) {
	defer i.wg.Done()

	_ = filepath.Walk(
		dir,
		func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if !fi.IsDir() {
				return nil
			} else if path.Base(p)[0] == '.' {
				return filepath.SkipDir
			}

			i.wg.Add(1)
			go i.index(bucket, p, fn)
			return nil
		},
	)
}

func (i *Index) index(bucket []byte, dir string, fn Indexer) {
	defer i.wg.Done()

	keys, err := fn(dir)

	if err == nil {
		err = i.db.Batch(func(tx *bolt.Tx) error {
			return writeKeys(tx.Bucket(bucket), keys, dir)
		})
	}

	if err != nil {
		i.err.Store(err)
	}
}

func (i *Index) gc() {
	_ = i.db.Update(func(tx *bolt.Tx) error {
		active := getActiveBucket(tx)

		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			if bytes.Equal(name, active) || bytes.Equal(name, metaBucketName) {
				return nil
			}

			return tx.DeleteBucket(name)
		})
	})
}

// WriteTo dumps a string representation of the database to w.
func (i *Index) WriteTo(w io.Writer) (int64, error) {
	var size int
	return int64(size), i.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			n, err := writeBucket(w, b, name, "")
			size += n
			return err
		})
	})
}

func writeBucket(w io.Writer, b *bolt.Bucket, name []byte, indent string) (int, error) {
	size, err := fmt.Fprintf(w, "%s* %s\n", indent, name)

	if err == nil {
		err = b.ForEach(func(k []byte, v []byte) error {
			var (
				n int
				e error
			)

			if v == nil {
				n, e = writeBucket(w, b.Bucket(k), k, indent+"  ")
			} else if len(v) == 0 {
				n, e = fmt.Fprintf(w, "%s  - '%s' (empty)\n", indent, k)
			} else {
				n, e = fmt.Fprintf(w, "%s  - '%s' : '%s'\n", indent, k, v)
			}

			size += n
			return e
		})
	}

	return size, err
}
