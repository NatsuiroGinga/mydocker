package bitmap

var bitValue = [32]int{
	0x00000001, 0x00000002, 0x00000004, 0x00000008,
	0x00000010, 0x00000020, 0x00000040, 0x00000080,
	0x00000100, 0x00000200, 0x00000400, 0x00000800,
	0x00001000, 0x00002000, 0x00004000, 0x00008000,
	0x00010000, 0x00020000, 0x00040000, 0x00080000,
	0x00100000, 0x00200000, 0x00400000, 0x00800000,
	0x01000000, 0x02000000, 0x04000000, 0x08000000,
	0x10000000, 0x20000000, 0x40000000, 0x80000000,
}

type BitMap struct {
	length  int64 // 总bit数
	bitsMap []int
}

func NewBitMap(length int64) *BitMap {
	/**
	 * 根据长度算出，所需数组大小
	 * 当 length%32=0 时大小等于
	 * = length/32
	 * 当 length%32>0 时大小等于
	 * = length/32+l
	 */
	length >>= 5
	if length&31 > 0 {
		length++
	}

	return &BitMap{
		length:  length,
		bitsMap: make([]int, length),
	}
}

