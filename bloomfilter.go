package inbloom

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
)

//BloomFilter is a space-efficient probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set.
type BloomFilter struct {
	vector         []byte
	baseHashFn     hash.Hash64
	numberOfHashes uint64
}

//Add an object to the set
func (filter *BloomFilter) Add(obj *[]byte) error {

	hashValues, err := filter.getHashVector(obj)

	for i := 0; i < len(hashValues); i++ {
		hashVal := hashValues[i]
		filter.vector[hashVal] = 1
	}
	return err
}

func (filter *BloomFilter) getHashVector(obj *[]byte) ([]uint64, error) {

	defer filter.baseHashFn.Reset()

	if filter == nil {
		fmt.Println("filter is null")
	}

	if filter.baseHashFn == nil {
		fmt.Println("hash function is null")
	}
	_, e1 := filter.baseHashFn.Write(*obj)

	if e1 != nil {
		empty := make([]uint64, 0)
		return empty, errors.New("failed to add object to filter")
	}

	seed1 := filter.baseHashFn.Sum64()
	// simulate output of two hash functions that will serve as base hashes for simulating the output of n hash functions
	upperBits := seed1 >> 32 << 32
	lowerBits := seed1 >> 32 << 32
	hashValues := make([]uint64, filter.numberOfHashes)
	for i := uint64(0); i < filter.numberOfHashes; i++ {

		h := (upperBits + lowerBits*i) % filter.numberOfHashes
		hashValues[i] = h
	}

	return hashValues, nil
}

//Test if an element is in the set
func (filter *BloomFilter) Test(obj *[]byte) (bool, error) {

	hashVales, e := filter.getHashVector(obj)

	for i := 0; i < len(hashVales); i++ {
		hashVal := hashVales[i]
		if filter.vector[hashVal] == 0 {
			return false, e
		}
	}

	return true, e
}

//NewFilter creates a new BloomFilter. p is the error rate
//and n is the estimated number of elements that will be handled by the filter
func NewFilter(p float64, n int64) BloomFilter {

	/*
		Given:

		n: how many items you expect to have in your filter (e.g. 216,553)
		p: your acceptable false positive rate {0..1} (e.g. 0.01 → 1%)
		we want to calculate:

		m: the number of bits needed in the bloom filter
		k: the number of hash functions we should apply
		The formulas:

		m = -n*ln(p) / (ln(2)^2) the number of bits
		k = m/n * ln(2) the number of hash functions

	*/

	m := -float64(n) * math.Log(p) / math.Pow(math.Ln2, 2)
	k := uint64(m / float64(n) * math.Ln2)

	return BloomFilter{vector: make([]byte, int64(m)),
		baseHashFn:     fnv.New64(),
		numberOfHashes: k}
}
