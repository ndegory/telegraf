package selfstat

import (
	"hash/fnv"
	"sort"
	"sync/atomic"
)

type stat struct {
	v           int64
	measurement string
	field       string
	tags        map[string]string
	key         uint64
}

func (s *stat) Incr(v int64) {
	atomic.AddInt64(&s.v, v)
}

func (s *stat) Set(v int64) {
	atomic.StoreInt64(&s.v, v)
}

func (s *stat) Get() int64 {
	return atomic.LoadInt64(&s.v)
}

func (s *stat) Name() string {
	return s.measurement
}

func (s *stat) FieldName() string {
	return s.field
}

// Tags returns a copy of the stat's tags.
// NOTE this allocates a new map every time it is called.
func (s *stat) Tags() map[string]string {
	m := make(map[string]string, len(s.tags))
	for k, v := range s.tags {
		m[k] = v
	}
	return m
}

func (s *stat) Key() uint64 {
	if s.key == 0 {
		h := fnv.New64a()
		h.Write([]byte(s.measurement))

		tmp := make([]string, len(s.tags)*2)
		i := 0
		for k, v := range s.tags {
			tmp[i] = k
			i++
			tmp[i] = v
			i++
		}
		sort.Strings(tmp)

		for _, s := range tmp {
			h.Write([]byte(s))
		}

		s.key = h.Sum64()
	}
	return s.key
}
