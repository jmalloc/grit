package source

import (
	"context"
	"sync"
)

// Repo is a reference to a remote repository.
type Repo struct {
	Name          string
	Description   string
	URL           string
	DefaultBranch string
}

// Resolver resolves a repository name to a single repository from a specific
// source.
type Resolver struct {
	Sources []Source
}

// ResolveMatch encapsulates a repository that is matched during a resolve
// operation.
type ResolveMatch struct {
	Source Source
	Repo   Repo
}

// SearchError encapsulates an error that occurred while searching a specific
// source during a resolve operation.
type SearchError struct {
	Source Source
	Error  error
}

// Selector is a function that chooses a single resolve match from multiple
// possibilities.
//
// It returns the desired match from the set of matches sent to the matches
// channel. The matches channel is closed when all repository sources have been
// searched.
type Selector func(
	ctx context.Context,
	matches <-chan ResolveMatch,
	errors <-chan SearchError,
) (m ResolveMatch, ok bool, _ error)

// Resolve searches repository sources for a single repository.
func (r *Resolver) Resolve(
	ctx context.Context,
	name string,
	selector Selector,
) (ResolveMatch, bool, error) {
	// Setup a context that we cancel to stop pending searches if a choice is
	// made before all searches are complete.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start searches in the background and wait for a choice to be made.
	matches := make(chan ResolveMatch)
	errors := make(chan SearchError)
	go r.searchAllSources(ctx, name, matches, errors)

	match, ok, err := selector(ctx, matches, errors)

	// Stop any pending searches, and wait for the matches channel to be closed
	// indicating that all search goroutines have finished.
	cancel()
	for range matches {
	}

	return match, ok, err
}

// searchAllSources searches for matching repositories from each source in
// parallel.
func (r *Resolver) searchAllSources(
	ctx context.Context,
	name string,
	matches chan<- ResolveMatch,
	errors chan<- SearchError,
) {
	var g sync.WaitGroup

	for _, source := range r.Sources {
		source := source // capture loop variable

		g.Add(1)
		go func() {
			defer g.Done()
			r.searchSingleSource(ctx, source, name, matches, errors)
		}()
	}

	g.Wait()
	close(matches)
}

// searchSingleSource searches a single source for matching repositories.
func (r *Resolver) searchSingleSource(
	ctx context.Context,
	source Source,
	name string,
	matches chan<- ResolveMatch,
	errors chan<- SearchError,
) {
	repos, err := source.Search(ctx, name)
	if err != nil {
		select {
		case <-ctx.Done():
		case errors <- SearchError{source, err}:
		}

		return
	}

	for _, repo := range repos {
		select {
		case <-ctx.Done():
			return
		case matches <- ResolveMatch{source, repo}:
		}
	}
}
