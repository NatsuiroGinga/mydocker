package bitmap

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type Hash32Func func(source string) int32

var Murmur3 = Hash32Func(func(source string) int32 {
	hasher := murmur3.New32()
	hasher.Write([]byte(source))
	return int32(hasher.Sum32() % math.MaxInt32)
})
