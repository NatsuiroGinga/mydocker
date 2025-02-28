package test

import (
	"encoding/json"
	"testing"
)

type A struct {
	F func(string) int
}

func Test_RemoveAll(t *testing.T) {
	a := &A{
		F: func(s string) int {
			return 0
		},
	}
	bytes, err := json.Marshal(a)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}
