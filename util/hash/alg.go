package hash

// Algorithm The algorithm for the hash
type Algorithm int

const (
	// Md5 The md5 algorithm
	Md5 Algorithm = iota
	// Sha1 The sha1 algorithm
	Sha1
	// Sha256 The sha256 algorithm
	Sha256
	// Sha3New256 The sha3-256 algorithm
	Sha3New256
)
