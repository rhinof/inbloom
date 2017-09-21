package inbloom

import (
	"errors"
	"hash"
	"hash/fnv"
	"math"
	"strings"
)

//ProbabilisticSet represents an abstraction of a Probabilistic
type ProbabilisticSet interface {
	Add(obj *[]byte) error
	Test(obj *[]byte) (bool, error)
}

//BloomFilter is a space-efficient probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set.
type BloomFilter struct {
	vector         []byte
	baseHashFn     hash.Hash64
	numberOfHashes uint64
}

//Add an object to the set
func (filter BloomFilter) Add(obj *[]byte) error {

	hashValues, err := filter.getHashVector(obj)

	for i := 0; i < len(hashValues); i++ {
		hashVal := hashValues[i]
		filter.vector[hashVal] = 1
	}
	return err
}

func (filter *BloomFilter) getHashVector(obj *[]byte) ([]uint64, error) {

	defer filter.baseHashFn.Reset()

	_, e1 := filter.baseHashFn.Write(*obj)

	if e1 != nil {
		empty := make([]uint64, 0)
		return empty, errors.New("failed to add object to filter")
	}

	// simulate output of two hash functions that will serve as base hashes for simulating the output of n hash functions
	// based on https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf

	// instead of generating 2 hash values we optimize and generate only 1 and splitting it's value
	// into 2 distinct values simulating the values of two separate hash functions
	seed1 := filter.baseHashFn.Sum64()

	upperBits := seed1 >> 32 << 32
	lowerBits := seed1 >> 32 << 32
	hashValues := make([]uint64, filter.numberOfHashes)

	for i := uint64(0); i < filter.numberOfHashes; i++ {

		h := (upperBits + lowerBits*i + uint64(i*i)) % filter.numberOfHashes
		hashValues[i] = h
	}

	return hashValues, nil
}

//Test if an element is in the set
func (filter BloomFilter) Test(obj *[]byte) (bool, error) {

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
func NewFilter(p float64, n int64) (filter ProbabilisticSet, e error) {

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

	defer func() {

		if r := recover(); r != nil {

			filter = nil
			err, ok := r.(error)

			if ok {
				msg := err.Error()
				if strings.Contains(msg, "makeslice: len out of range") {
					e = errors.New("failed to create BloomFilter. either reduce error rate or maximum estimated elements in filter ")
				} else {
					e = err
				}

			} else {
				e = errors.New("failed to create BloomFilter")
			}

		}
	}()

	m := -float64(n) * math.Log(p) / math.Pow(math.Ln2, 2)
	k := uint64(m / float64(n) * math.Ln2)
	sliceLength := int64(m)

	bf := BloomFilter{vector: make([]byte, sliceLength),
		baseHashFn:     fnv.New64(),
		numberOfHashes: k}

	return bf, nil

}
