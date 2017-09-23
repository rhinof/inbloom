package inbloom

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestAddToFilter(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()
	filter, _ := NewFilter(0.1, 10000)
	err := filter.Add(&data)
	if err != nil {
		t.Fail()
	}

}

func TestNotInFilterReturnsFalse(t *testing.T) {

	data := bytes.NewBufferString("rhinof is on the moo").Bytes()
	data1 := bytes.NewBufferString("shama is the king").Bytes()
	filter, _ := NewFilter(0.1, 10000)
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

	filter, e := NewFilter(0.1, 1000000000)

	if e != nil {
		t.Errorf("%s", e)
		t.FailNow()
	}

	err := filter.Add(&data)

	if err != nil {

		t.Fail()
	}

	result, _ := filter.Test(&data)

	if result == false {
		t.Fail()
	}

}

func TestPoFPInRange(t *testing.T) {

	n := 100000
	p := 0.01
	filter, _ := NewFilter(p, n)
	i := 0

	for filter.PoFP() <= p {
		data := make([]byte, 30)
		rand.Read(data)
		filter.Add(&data)
		i++

	}

	if float64(n)*float64(0.9) > float64(i) {
		t.Fatal("filter error rate wasn't in acceptable range")
	}

}

func BenchmarkFilter(b *testing.B) {

	filter, _ := NewFilter(0.01, 100000)

	for i := 0; i < 100000; i++ {

		data := bytes.NewBufferString("rhinof is on the moo" + string(i)).Bytes()

		err := filter.Add(&data)
		if err != nil {
			b.Fail()
		}

		result, _ := filter.Test(&data)

		if result == false {
			b.Fail()
		}
	}
}
