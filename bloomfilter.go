package inbloom

import (
	"bytes"
	"encoding/gob"
	"errors"
	"hash"
	"hash/fnv"
	"io"
	"math"
)

type BloomFilter struct {
	vector         []byte
	baseHashFn     hash.Hash64
	numberOfHashes uint64
}

func (self *BloomFilter) Add(obj *[]byte) error {

	hashValues, err := self.getHashVector(obj)

	for i := 0; i < len(hashValues); i++ {
		hashVal := hashValues[i]
		self.vector[hashVal] = 1
	}
	return err
}

func (self *BloomFilter) getHashVector(obj *[]byte) ([]uint64, error) {

	defer self.baseHashFn.Reset()

	_, e1 := self.baseHashFn.Write(*obj)

	if e1 != nil {
		empty := make([]uint64, 0)
		return empty, errors.New("failed to add object to filter")
	}

	seed1 := self.baseHashFn.Sum64()
	// simulate output of two hash functions that will serve as base hashes for simulating the output of n hash functions
	upperBits := seed1 >> 32 << 32
	lowerBits := seed1 >> 32 << 32
	hashValues := make([]uint64, self.numberOfHashes)
	for i := uint64(0); i < self.numberOfHashes; i++ {

		h := (upperBits + lowerBits*i) % self.numberOfHashes
		hashValues[i] = h
	}

	return hashValues, nil
}

func (self *BloomFilter) Test(obj *[]byte) (bool, error) {

	hashVales, e := self.getHashVector(obj)

	for i := 0; i < len(hashVales); i++ {
		hashVal := hashVales[i]
		if self.vector[hashVal] == 0 {
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

func Inflate(reader io.Reader) *BloomFilter {

	var filter BloomFilter
	dec := gob.NewDecoder(reader)
	dec.Decode(&filter)
	return &filter
}

func Deflate(filter *BloomFilter) []byte {

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(filter)

	return buffer.Bytes()
}
