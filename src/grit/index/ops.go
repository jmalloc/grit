package index

import "github.com/boltdb/bolt"

var slugBucket = []byte("slugs")
var dirBucket = []byte("dirs")

func update(tx *bolt.Tx, dir string, slugs set) error {
	if err := updateDir(tx, dir, dirRecord{
		Slugs: slugs,
	}); err != nil {
		return err
	}

	for s := range slugs {
		if err := mergeSlug(tx, s, slugRecord{
			Dirs: newSet(dir),
		}); err != nil {
			return err
		}
	}

	return nil
}

func updateDir(tx *bolt.Tx, dir string, r dirRecord) error {
	bucket, err := tx.CreateBucketIfNotExists(dirBucket)
	if err != nil {
		return err
	}

	k := []byte(dir)
	return bucket.Put(k, pack(r))
}

func mergeSlug(tx *bolt.Tx, slug string, r slugRecord) error {
	bucket, err := tx.CreateBucketIfNotExists(slugBucket)
	if err != nil {
		return err
	}

	k := []byte(slug)
	rec := unpackSlugRecord(bucket.Get(k))
	rec.Dirs.Merge(r.Dirs)

	return bucket.Put(k, pack(rec))
}

func remove(tx *bolt.Tx, dir string) error {
	bucket := tx.Bucket(dirBucket)
	if bucket == nil {
		return nil
	}

	k := []byte(dir)
	rec := unpackDirRecord(bucket.Get(k))

	for s := range rec.Slugs {
		if err := removeDirFromSlug(tx, dir, s); err != nil {
			return err
		}
	}

	return bucket.Delete(k)
}

func removeDirFromSlug(tx *bolt.Tx, dir, slug string) error {
	bucket := tx.Bucket(slugBucket)
	if bucket == nil {
		return nil
	}

	k := []byte(slug)
	rec := unpackSlugRecord(bucket.Get(k))
	rec.Dirs.Remove(dir)

	if len(rec.Dirs) == 0 {
		return bucket.Delete(k)
	}

	return bucket.Put(k, pack(rec))
}
