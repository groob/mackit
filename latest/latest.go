package latest

import (
	"regexp"
	"sort"
	"strconv"
)

// Version returns the latest macOS build number given a list of them.
func Version(versions ...string) string {
	if len(versions) == 0 {
		return ""
	}

	v := Sorted(versions...)
	return v[0]
}

// Sorted sorts a list of macOS build numbers with the latest build number first.
func Sorted(versions ...string) []string {
	v := byVersion(versions)
	sort.Sort(v)
	return v
}

// ReverseSorted sorts a list of macOS build numbers from oldest to newest.
func ReverseSorted(versions ...string) []string {
	v := byVersion(versions)
	sort.Sort(sort.Reverse(v))
	return v
}

var r = regexp.MustCompile(`(?P<major>[\d]+)(?P<min>[A-Z])(?P<patch>[\d]+)(?P<fix>[a-z]?)`)

type byVersion []string

func (a byVersion) Len() int      { return len(a) }
func (a byVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byVersion) Less(i, j int) bool {
	return a[i] == version(a[i], a[j])
}

// return the newer of the two
func version(a, b string) string {
	if a == b {
		return a
	}

	aa := r.FindStringSubmatch(a)
	bb := r.FindStringSubmatch(b)

	maa, err := strconv.Atoi(aa[1])
	if err != nil {
		panic(err)
	}

	mbb, err := strconv.Atoi(bb[1])
	if err != nil {
		panic(err)
	}

	major := maa > mbb
	if major {
		return a
	} else if maa < mbb {
		return b
	}

	minor := aa[2] > bb[2]
	if minor {
		return a
	} else if aa[2] < bb[2] {
		return b
	}

	paa, err := strconv.Atoi(aa[3])
	if err != nil {
		panic(err)
	}

	pbb, err := strconv.Atoi(bb[3])
	if err != nil {
		panic(err)
	}

	patch := paa > pbb
	if patch {
		return a
	} else if paa < pbb {
		return b
	}

	if len(aa) > len(bb) {
		// is a re-build with `a-z` at the end.
		return a
	}

	if len(aa) == len(bb) && len(aa) == 5 {
		// both have a re-build letter, compare.
		fixup := aa[4] > bb[4]
		if fixup {
			return a
		}
	}

	return b
}
