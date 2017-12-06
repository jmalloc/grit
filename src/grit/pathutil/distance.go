package pathutil

import (
	"math"
	"path/filepath"
	"strings"
)

// Distance returns the "distance" from base to target. The lower the
// distance the "closer" target is to "base".
func Distance(base, target string) (dist uint32) {
	if base == target {
		return 0
	}

	rel, err := filepath.Rel(base, target)
	if err != nil {
		return math.MaxUint32
	}

	// count the number of path separators to get the distance
	dist = uint32(strings.Count(
		rel,
		string(filepath.Separator),
	))

	// increase the distance for targets that are not a subfolder of base
	if strings.HasPrefix("..", rel) {
		dist += math.MaxUint32 / 2
	}

	return
}
