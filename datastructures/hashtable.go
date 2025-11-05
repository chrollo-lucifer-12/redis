package datastructures

type HashSet struct {
	buckets     [][]string
	numElements int
}

const loadFactorLimit = 0.75

func NewHashSet() *HashSet {
	return &HashSet{
		buckets: make([][]string, 8),
	}
}

func hashString(s string) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(s); i++ {
		hash = ((hash << 5) + hash) + uint64(s[i])
	}
	return hash
}

func (s *HashSet) getIndex(key string) int {
	return int(hashString(key) % uint64(len(s.buckets)))
}

func (s *HashSet) Insert(key string) {
	index := s.getIndex(key)

	for _, val := range s.buckets[index] {
		if val == key {
			return
		}
	}

	s.buckets[index] = append(s.buckets[index], key)
	s.numElements++

	if float64(s.numElements)/float64(len(s.buckets)) > loadFactorLimit {
		s.rehash()
	}
}

func (s *HashSet) Contains(key string) bool {
	index := s.getIndex(key)

	for _, val := range s.buckets[index] {
		if val == key {
			return true
		}
	}
	return false
}

func (s *HashSet) rehash() {
	newBucketCount := len(s.buckets) * 2
	newBuckets := make([][]string, newBucketCount)

	for _, bucket := range s.buckets {
		for _, val := range bucket {
			newIndex := int(hashString(val) % uint64(newBucketCount))
			newBuckets[newIndex] = append(newBuckets[newIndex], val)
		}
	}

	s.buckets = newBuckets
}

func (s *HashSet) Size() int {
	return s.numElements
}

func (s *HashSet) Erase(key string) bool {
	index := s.getIndex(key)
	bucket := s.buckets[index]
	for i, val := range bucket {
		if val == key {
			s.buckets[index] = append(bucket[:i], bucket[i+1:]...)
			s.numElements--
			return true
		}
	}
	return false
}
func (s *HashSet) Elements() []string {
	result := make([]string, 0, s.numElements)
	for _, bucket := range s.buckets {
		for _, val := range bucket {
			result = append(result, val)
		}
	}
	return result
}
