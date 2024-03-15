package grit

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport"
)

func TestReplaceSlug_GitUrl(t *testing.T) {
	ep, err := transport.NewEndpoint("git@github.com:owner-a/repo-a.git")
	if err != nil {
		t.Fatal(err)
	}

	replaced := ReplaceSlug(ep, "owner-b/repo-b")

	if replaced.String() != "ssh://git@github.com/owner-b/repo-b.git" {
		t.Errorf("expected ssh://git@github.com/owner-b/repo-b.git, got %s", replaced.String())
	}
}

func TestReplaceSlug_SshUrl(t *testing.T) {
	ep, err := transport.NewEndpoint("ssh://git@github.com/owner-a/repo-a.git")
	if err != nil {
		t.Fatal(err)
	}

	replaced := ReplaceSlug(ep, "owner-b/repo-b")

	if replaced.String() != "ssh://git@github.com/owner-b/repo-b.git" {
		t.Errorf("expected ssh://git@github.com/owner-b/repo-b.git, got %s", replaced.String())
	}
}

func TestReplaceSlug_HttpUrl(t *testing.T) {
	ep, err := transport.NewEndpoint("https://github.com/owner-a/repo-a.git")
	if err != nil {
		t.Fatal(err)
	}

	replaced := ReplaceSlug(ep, "owner-b/repo-b")

	if replaced.String() != "https://github.com/owner-b/repo-b.git" {
		t.Errorf("expected https://github.com/owner-b/repo-b.git, got %s", replaced.String())
	}
}
