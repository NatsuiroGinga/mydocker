package test

import (
	"os"
	"testing"
)

func Test_RemoveAll(t *testing.T) {
	os.RemoveAll("../test-remove")
}
