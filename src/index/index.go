package index

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/boltdb/bolt"
	"github.com/jmalloc/grit/src/provider"
)

// Index is an index of repository locations.
type Index struct {
	db        *bolt.DB
	providers []*provider.Provider
}

// Open opens the index database at path f.
func Open(f string, p []*provider.Provider) (*Index, error) {
	db, err := bolt.Open(f, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Index{db, p}, nil
}

// Rebuild the entire index.
func (i *Index) Rebuild(paths ...string) error {
	defer i.gc()

	for _, p := range i.providers {
		paths = append(paths, p.BasePath)
	}

	bucket, err := i.newBucket()
	if err != nil {
		return err
	}

	pending := 0
	errors := make(chan error)

	for _, p := range paths {
		_ = filepath.Walk(p, func(dir string, info os.FileInfo, err error) error {
			if _, err := os.Stat(path.Join(dir, ".git")); err != nil {
				return nil
			}

			pending++
			go func() {
				errors <- i.indexClone(bucket, dir)
			}()

			return filepath.SkipDir
		})
	}

	for e := range errors {
		if e != nil {
			err = e
		}
		pending--
		if pending == 0 {
			close(errors)
		}
	}

	if err == nil {
		return i.activateBucket(bucket)
	}

	return err
}

func (i *Index) indexClone(bucket []byte, dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	slugs, err := i.slugs(repo)
	if err != nil {
		return err
	} else if len(slugs) == 0 {
		return nil
	}

	return i.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		for _, slug := range slugs {
			if err := b.Put([]byte(slug), []byte(dir)); err != nil {
				return err
			}
		}

		return nil
	})
}

// Close closes the index.
func (i *Index) Close() {
	_ = i.db.Close()
}

// WriteTo dumps a string representation of the database to w.
func (i *Index) WriteTo(w io.Writer) (n int64, err error) {
	err = i.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			s, err := fmt.Fprintf(w, "- %s\n", name)
			n += int64(s)
			if err != nil {
				return err
			}
			return b.ForEach(func(k []byte, v []byte) error {
				s, err := fmt.Fprintf(w, "  - %s = %s\n", k, v)
				n += int64(s)
				return err
			})
		})
	})

	return
}

func (i *Index) newBucket() (name []byte, err error) {
	err = i.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(metaBucket)
		if err != nil {
			return err
		}

		seq, _ := b.NextSequence()
		name = []byte(fmt.Sprintf("repos-%d", seq))

		_, err = tx.CreateBucket(name)
		return err
	})

	return
}

func (i *Index) activateBucket(name []byte) error {
	return i.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(metaBucket)
		if err != nil {
			return err
		}

		return b.Put(activeBucketKey, name)
	})
}

func (i *Index) gc() {
	_ = i.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(metaBucket)
		if b == nil {
			return nil
		}

		active := b.Get(activeBucketKey)

		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if bytes.Equal(name, active) || bytes.Equal(name, metaBucket) {
				return nil
			}

			return tx.DeleteBucket(name)
		})
	})
}

func (i *Index) slugs(r *git.Repository) (slugs []string, err error) {
	var s []string
	for _, p := range i.providers {
		s, err = p.Driver.Slugs(r)
		if err != nil {
			return
		}

		slugs = append(slugs, s...)
	}

	return
}

var (
	metaBucket      = []byte("meta")
	activeBucketKey = []byte("active")
)
