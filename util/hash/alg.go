package hash

type Algorithm int

const (
	Md5 Algorithm = iota
	Sha1
	Sha256
	Sha3New256
)
