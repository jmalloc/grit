package index

type set map[string]struct{}

func newSet(keys ...string) set {
	s := set{}
	s.Add(keys...)
	return s
}

func (s set) Add(keys ...string) {
	for _, k := range keys {
		s[k] = struct{}{}
	}
}

func (s set) Remove(keys ...string) {
	for _, k := range keys {
		delete(s, k)
	}
}

func (s set) Keys() []string {
	var keys []string
	for str := range s {
		keys = append(keys, str)
	}
	return keys
}

func (s set) Merge(o set) {
	for k := range o {
		s[k] = struct{}{}
	}
}
