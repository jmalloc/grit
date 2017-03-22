package repo

import (
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
)

// Exists returns true if url is a Git repository.
func Exists(url string) (ok bool, err error) {
	endpoint, err := transport.NewEndpoint(url)
	if err != nil {
		return
	}

	cli, err := client.NewClient(endpoint)
	if err != nil {
		return
	}

	sess, err := cli.NewUploadPackSession(endpoint, nil)
	if err != nil {
		return
	}
	defer sess.Close()

	_, err = sess.AdvertisedReferences()

	if err == transport.ErrRepositoryNotFound {
		err = nil
	} else if err == nil {
		ok = true
	}

	return
}
