package inbloom

import (
	"bytes"
	"testing"
)

func TestAddToFilter(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()
	filter := NewFilter(0.1, 10000)
	err := filter.Add(&data)
	if err != nil {
		t.Fail()
	}

}

func TestNotInFilterReturnsFalse(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()
	data1 := bytes.NewBufferString("shama is the king").Bytes()
	filter := NewFilter(0.1, 10000)
	err := filter.Add(&data1)
	if err != nil {
		t.Fail()
	}

	result, _ := filter.Test(&data)

	if result == true {
		t.Fail()
	}

}

func TestDataIsInFilterReturnsTrue(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()

	filter := NewFilter(0.1, 10000)
	err := filter.Add(&data)
	if err != nil {
		t.Fail()
	}

	result, _ := filter.Test(&data)

	if result == false {
		t.Fail()
	}

}

func Testserializing(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()
	f := bytes.NewBufferString("not if filter").Bytes()

	filter := NewFilter(0.1, 10000)
	filter.Add(&data)
	deflated := Deflate(&filter)

	buffer := bytes.NewBuffer(deflated)
	inflatedFilter := Inflate(buffer)
	result, _ := inflatedFilter.Test(&data)
	result1, _ := inflatedFilter.Test(&f)

	if result1 == true {
		t.Fail()
	}
	if result == false {
		t.Fail()
	}

}
