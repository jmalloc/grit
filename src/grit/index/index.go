package index

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"sync/atomic"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/boltdb/bolt"
)

// EndpointFilter is a function that returns true if the given endpoint should
// be included for indexing when scanning.
type EndpointFilter func(ep transport.Endpoint) bool

// Index is an index of repository locations.
type Index struct {
	db *bolt.DB

	sm  sync.Mutex // protects all state below
	wm  sync.Mutex // protects w from concurrent writes
	w   io.Writer
	f   EndpointFilter
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

// Add a clone path to the index.
func (i *Index) Add(dir string) error {
	slugs, err := slugsFromClone(dir, nil)
	if err != nil {
		return err
	}

	return i.db.Update(func(tx *bolt.Tx) error {
		return update(tx, dir, slugs)
	})
}

// Remove removes a clone path from the index.
func (i *Index) Remove(dir string) error {
	return i.db.Update(func(tx *bolt.Tx) error {
		return remove(tx, dir)
	})
}

// Has returns true if dir is in the index.
func (i *Index) Has(dir string) (ok bool) {
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dirBucket)
		if bucket != nil {
			k := []byte(dir)
			ok = bucket.Get(k) != nil
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return
}

// Find returns a list of paths matching the given slug.
func (i *Index) Find(slug string) []string {
	var rec slugRecord
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(slugBucket)
		if bucket != nil {
			k := []byte(slug)
			rec = unpackSlugRecord(bucket.Get(k))
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return rec.Dirs.Keys()
}

// FindByDir returns a list of slugs that match dir. dir may be a
// sub-directory of the clone. d is the indexed directory, i.e. the clone root.
func (i *Index) FindByDir(dir string) (s []string, d string) {
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dirBucket)
		if bucket != nil {
			d = dir
			for d != "." && d != "/" {
				k := []byte(d)
				v := bucket.Get(k)

				if v != nil {
					rec := unpackDirRecord(v)
					s = rec.Slugs.Keys()
					return nil
				}

				d = path.Dir(d)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return
}

// ListSlugs returns the slugs that begin with p, which may be empty.
func (i *Index) ListSlugs(p string) []string {
	var slugs []string

	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(slugBucket)
		if bucket == nil {
			return nil
		}

		pre := []byte(p)
		return bucket.ForEach(func(slug []byte, _ []byte) error {
			if bytes.HasPrefix(slug, pre) {
				slugs = append(slugs, string(slug))
			}
			return nil
		})
	})

	if err != nil {
		panic(err)
	}

	return slugs
}

// ListClones returns a list of all clone directories.
func (i *Index) ListClones() []string {
	var dirs []string

	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dirBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(dir []byte, _ []byte) error {
			dirs = append(dirs, string(dir))
			return nil
		})
	})

	if err != nil {
		panic(err)
	}

	return dirs
}

// Prune removes directories that no longer exist.
func (i *Index) Prune(w io.Writer) error {
	return i.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dirBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(dir []byte, buf []byte) error {
			s := string(dir)
			if isDir(s) {
				return nil
			}

			if w != nil {
				fmt.Fprintln(w, s)
			}

			return remove(tx, s)
		})
	})
}

// Scan recursively indexes dirs.
func (i *Index) Scan(
	w io.Writer,
	f EndpointFilter,
	dirs ...string,
) error {
	i.sm.Lock()
	defer i.sm.Unlock()

	i.w = w
	i.f = f

	for _, d := range dirs {
		i.wg.Add(1)
		go i.scan(d)
	}

	i.wg.Wait()
	err, _ := i.err.Load().(error)
	return err
}

func (i *Index) scan(dir string) {
	var err error
	defer func() {
		i.wg.Done()
		if err != nil {
			i.err.Store(err)
		}
	}()

	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}

	err = filepath.Walk(dir, i.walk)
}

func (i *Index) walk(dir string, info os.FileInfo, err error) error {
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return nil
	}

	if path.Base(dir)[0] == '.' {
		return filepath.SkipDir // don't index hidden directories
	}

	if isGitDir(dir) {
		i.wg.Add(1)
		go i.batch(dir)
		return filepath.SkipDir // skip sub-directories of git clones
	}

	return nil
}

func (i *Index) batch(dir string) {
	var err error
	defer func() {
		i.wg.Done()
		if err != nil {
			i.err.Store(err)
		}
	}()

	slugs, err := slugsFromClone(dir, i.f)
	if err != nil || len(slugs) == 0 {
		return
	}

	i.wm.Lock()
	fmt.Fprintln(i.w, dir)
	i.wm.Unlock()

	err = i.db.Batch(func(tx *bolt.Tx) error {
		return update(tx, dir, slugs)
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

func isDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func isGitDir(dir string) bool {
	return isDir(path.Join(dir, ".git"))
}
