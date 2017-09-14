package inbloom

import (
	"bytes"
	"testing"
)

func TestAddToFilter(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo")
	filter := NewFilter(0.1, 10000)
	err := filter.Add(data.Bytes())
	if err != nil {
		t.Fail()
	}

}

func TestNotInFilterReturnsFalse(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo")
	data1 := bytes.NewBufferString("shama is the king")
	filter := NewFilter(0.1, 10000)
	err := filter.Add(data1.Bytes())
	if err != nil {
		t.Fail()
	}

	result, _ := filter.Test(data.Bytes())

	if result == true {
		t.Fail()
	}

}

func TestDataIsInFilterReturnsTrue(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo")

	filter := NewFilter(0.1, 10000)
	err := filter.Add(data.Bytes())
	if err != nil {
		t.Fail()
	}

	result, _ := filter.Test(data.Bytes())

	if result == false {
		t.Fail()
	}

}

func Testserializing(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo")

	filter := NewFilter(0.1, 10000)
	filter.Add(data.Bytes())
	deflated := Deflate(&filter)

	buffer := bytes.NewBuffer(deflated)
	inflatedFilter := Inflate(buffer)
	result, _ := inflatedFilter.Test(data.Bytes())

	if result == false {
		t.Fail()
	}

}
