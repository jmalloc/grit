package grit

import (
	"bytes"
	"html/template"
	"os"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
)

// EndpointTemplate is template for a Git repository URL.
type EndpointTemplate string

// Endpoint represents a Git clone endpoint, resolved from an EndpointTemplate.
type Endpoint struct {
	// The actual URL used to clone the repository.
	// This string will match the URL template from the configuration as closely
	// as possible.
	Actual string

	// The substituted and normalized endpoint template. SCP-style Git URLs
	// are converted to ssh:// URLs.
	Normalized transport.Endpoint
}

// Validate returns an error if the template is invalid.
func (t EndpointTemplate) Validate() error {
	_, err := t.virtualEndpoint()
	return err
}

// IsMatch returns true if may have been derived from the endpoint template.
func (t EndpointTemplate) IsMatch(e transport.Endpoint) bool {
	ep, err := t.virtualEndpoint()
	if err != nil {
		return false
	}

	// TODO: match slug heuristically
	return e.Scheme == ep.Scheme && e.Host == ep.Host
}

// virtualEndpoint returns a Git endpoint from the template as though we had
// a slug to resolve.
func (t EndpointTemplate) virtualEndpoint() (transport.Endpoint, error) {
	ep, err := t.Resolve("__virtual__")
	return ep.Normalized, err
}

// Resolve returns a URL from the template.
func (t EndpointTemplate) Resolve(slug string) (ep Endpoint, err error) {
	ep.Actual, err = t.replace(slug)

	if err == nil {
		ep.Normalized, err = transport.NewEndpoint(ep.Actual)
	}

	return
}

func (t EndpointTemplate) replace(slug string) (u string, err error) {
	funcs := map[string]interface{}{
		"slug": func() string { return slug },
		"env":  os.Getenv,
	}

	tmpl, err := template.
		New("url").
		Funcs(funcs).
		Parse(string(t))

	if err == nil {
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, nil)
		u = buf.String()
	}

	return
}

// EndpointExists returns true if url is a Git repository.
func EndpointExists(ep Endpoint) (ok bool, err error) {
	cli, err := client.NewClient(ep.Normalized)
	if err != nil {
		return
	}

	sess, err := cli.NewUploadPackSession(ep.Normalized, nil)
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

// EndpointToDir returns the absolute path for a clone of a repository.
func EndpointToDir(base string, ep transport.Endpoint) string {
	p := strings.TrimSuffix(ep.Path, path.Ext(ep.Path))

	return path.Join(base, ep.Host+p)
}
