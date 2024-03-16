package bf

const (
	fnvPrime       uint64 = 1099511628211
	fnvOffsetBasis uint64 = 14695981039346656037
)

func fnv1Hash(data []byte) uint64 {
	hash := fnvOffsetBasis

	for _, databyte := range data {
		hash = hash * fnvPrime
		hash = hash ^ uint64(databyte)
	}

	return hash
}
func fnv1aHash(data []byte) uint64 {
	hash := fnvOffsetBasis

	for _, databyte := range data {
		hash = hash ^ uint64(databyte)
		hash = hash * fnvPrime
	}

	return hash
}
