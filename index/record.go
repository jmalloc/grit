package index

import "encoding/json"

func pack(v interface{}) []byte {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return buf
}

type dirRecord struct {
	Slugs set
}

func unpackDirRecord(buf []byte) (rec dirRecord) {
	rec.Slugs = set{}

	if buf != nil {
		if err := json.Unmarshal(buf, &rec); err != nil {
			panic(err)
		}
	}

	return
}

type slugRecord struct {
	Dirs set
}

func unpackSlugRecord(buf []byte) (rec slugRecord) {
	rec.Dirs = set{}

	if buf != nil {
		if err := json.Unmarshal(buf, &rec); err != nil {
			panic(err)
		}
	}

	return
}
