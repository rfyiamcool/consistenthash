package consistenhash

import (
	"hash/crc32"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

var (
	lock sync.RWMutex
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int   // Vnodes
	keys     []int // Sorted
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *Map) Add(keys ...string) {
	lock.Lock()
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hashCode := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hashCode)
			m.hashMap[hashCode] = key
		}
	}
	sort.Ints(m.keys)
	lock.Unlock()
}

func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// optimize
	idx := sort.Search(
		len(m.keys),
		func(i int) bool { return m.keys[i] >= hash },
	)

	if idx == len(m.keys) {
		idx = 0
	}

	lock.RLock()
	data := m.hashMap[m.keys[idx]]
	lock.RUnlock()

	return data
}

func HashToInt(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return uint32(h.Sum32())
}
