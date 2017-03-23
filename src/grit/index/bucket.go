package index

import (
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	metaBucketName  = []byte("meta")
	activeBucketKey = []byte("active")
)

func readMetaBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(metaBucketName)
}

func writeMetaBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	return tx.CreateBucketIfNotExists(metaBucketName)
}

func readActiveBucket(tx *bolt.Tx) *bolt.Bucket {
	active := getActiveBucket(tx)
	if active == nil {
		return nil
	}

	return tx.Bucket(active)
}

func writeActiveBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	meta, err := writeMetaBucket(tx)
	if err != nil {
		return nil, err
	}

	active := meta.Get(activeBucketKey)
	if active == nil {
		active, err = generateBucketName(meta)
		if err != nil {
			return nil, err
		}
	}

	return tx.CreateBucketIfNotExists(active)
}

func writeNextBucket(tx *bolt.Tx) ([]byte, *bolt.Bucket, error) {
	meta, err := writeMetaBucket(tx)
	if err != nil {
		return nil, nil, err
	}

	name, err := generateBucketName(meta)
	if err != nil {
		return nil, nil, err
	}

	bucket, err := tx.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, nil, err
	}

	return name, bucket, nil
}

func setActiveBucket(tx *bolt.Tx, name []byte) (err error) {
	meta, err := writeMetaBucket(tx)

	if err == nil {
		err = meta.Put(activeBucketKey, name)
	}

	return
}

func getActiveBucket(tx *bolt.Tx) []byte {
	meta := readMetaBucket(tx)
	if meta == nil {
		return nil
	}

	return meta.Get(activeBucketKey)
}

func writeKeys(bucket *bolt.Bucket, keys []string, dir string) error {
	dirBytes := []byte(dir)

	for _, key := range keys {
		keyBytes := []byte(key)
		sub, err := bucket.CreateBucketIfNotExists(keyBytes)
		if err != nil {
			return err
		}

		if err := sub.Put(dirBytes, []byte{}); err != nil {
			return err
		}
	}

	return nil
}

func generateBucketName(meta *bolt.Bucket) (name []byte, err error) {
	seq, err := meta.NextSequence()
	if err == nil {
		name = []byte(fmt.Sprintf("v%d", seq))
	}
	return
}
