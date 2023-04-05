package singleflight

import "testing"

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return "zz", nil
	})

	if v != "zz" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}
