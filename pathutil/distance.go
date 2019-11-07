package pathutil

import (
	"math"
	"path/filepath"
	"sort"
	"strings"
)

// Distance is a numerical distance between two paths
type Distance uint32

const (
	// PreferBase is the zero distance. When passed to SortByDistance(), the
	// base path is sorted to the top.
	PreferBase Distance = 0

	// PreferChildren is a distance midway between children of base paths, and
	// non-children. When passed to SortByDistance(), the base path is sorted
	// below its own children, but above non-children.
	PreferChildren Distance = math.MaxUint32 / 2

	// PreferOther is the maximimum distance. When passed to SortByDistance(),
	// the base path is sorted to the bottom.
	PreferOther Distance = math.MaxUint32
)

// DistanceBetween returns the "distance" from base to target. The more "hops"\
// required to move from base to target, the larger the distance. Targets that
// are children of base are "closer" than targets that are not. If target is
// equal to base, the distance is baseDist.
func DistanceBetween(base, target string, baseDist Distance) (dist Distance) {
	if base == target {
		return baseDist
	}

	rel, err := filepath.Rel(base, target)
	if err != nil {
		return math.MaxUint32
	}

	// count the number of path separators to get the distance
	dist = Distance(strings.Count(
		rel,
		string(filepath.Separator),
	))

	// increase the distance for targets that are not a subfolder of base
	if strings.HasPrefix(rel, "..") {
		dist += 1 + PreferChildren
	}

	return
}

// SortByDistance sorts paths according to their "distance" from base, according
// to Distance(base, path). Any two paths with the same distance from base are
// further sorted by their name.
func SortByDistance(base string, paths []string, baseDist Distance) {
	type elem struct {
		dist Distance
		path string
	}

	elems := make([]elem, len(paths))
	for i, p := range paths {
		elems[i] = elem{
			DistanceBetween(base, p, baseDist),
			p,
		}
	}

	sort.Slice(elems, func(i, j int) bool {
		ei, ej := elems[i], elems[j]

		if ei.dist == ej.dist {
			return strings.Compare(ei.path, ej.path) < 0
		}

		return ei.dist < ej.dist
	})

	for idx, elem := range elems {
		paths[idx] = elem.path
	}
}
