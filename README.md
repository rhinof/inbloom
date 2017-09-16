# inbloom
A fast bloom filter written in Go.

A Bloom filter is a space-efficient probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set.

https://en.wikipedia.org/wiki/Bloom_filter

### Usage

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/rhinof/inbloom"
)

func main() {

	//Probability of false positives
	fpRate := 0.1
	// Estimated set size
	expectedSetSize := int64(100000)
	data := bytes.NewBufferString("rhinof").Bytes()

	//Create the filter
	filter := inbloom.NewFilter(fpRate, expectedSetSize)

	if err := filter.Add(&data); err != nil {
		fmt.Print("couldn't add element")
	}

	result, _ := filter.Test(&data)

	fmt.Printf("test returned %t", result)
}
```

