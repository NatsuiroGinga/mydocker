package bitmap

import (
	"strconv"
)

type BloomFilter interface {
	Exist(val string) bool // 判断值是否在bitmap中
	Set(val string)        // 追加元素到bitmap中
}

// LocalBloomService 本地 bloom filter
type LocalBloomService struct {
	M      int32      // bitmap的长度，由用户输入
	K      int32      // hash函数的个数，由bloom filter统计
	N      int32      // bloom filter中的元素个数
	Bitmap []int32    // 位图，类型为[]int, 其中使用到每个int的32位，因此有[]int长度为m/32,构造时为了避免除不尽的问题，切片长度额外增加1
	hasher Hash32Func // hash函数编码模块
}

func NewLocalBloomService(m, k int32, hasher Hash32Func) *LocalBloomService {
	return &LocalBloomService{
		M:      m,
		K:      k,
		Bitmap: make([]int32, m/32+1),
		hasher: hasher,
	}
}

/*
由于[]int中,每个int元素使用32个bit位,因此对于每个offset,对应在[]int中的index位置为offset>>5,即offset/32

offset 在一个int元素中的位置对应为 offset &31,即 offset% 32

倘若有任意一个bit位标识为0,都说明元素val在布隆过滤器中一定不存在

倘若所有bit位标识都为1,则说明元素val在布隆过滤器中很有可能存在
*/
func (bloom *LocalBloomService) Exist(val string) bool {
	if bloom.hasher == nil {
		bloom.hasher = Murmur3
	}
	for _, offset := range bloom.getKEncrypted(val) {
		index := offset >> 5     // 等价于 / 32
		bitOffset := offset & 31 // 等价于 % 32

		if bloom.Bitmap[index]&(1<<bitOffset) == 0 {
			return false
		}
	}
	return true
}

/*
getKEncrypted 获取一个元素val对应k个bit位偏移量offset的实现如下:

1.首次映射时,以元素val作为输入,获取murmur3映射得到的hash值

2.接下来每次以上一轮的hash值作为输入,获取murmur3映射得到新一轮hash值

3.凑齐k个hash值后返回结果
*/
func (bloom *LocalBloomService) getKEncrypted(val string) []int32 {
	encrypteds := make([]int32, 0, bloom.K)
	origin := val

	for i := range bloom.K {
		encrypted := bloom.hasher(origin) & (bloom.M - 1)
		encrypteds = append(encrypteds, encrypted)
		if i == bloom.K-1 {
			break
		}
		origin = strconv.Itoa(int(encrypted))
	}

	return encrypteds
}

/*
下面是追加元素进入布隆过滤器的流程:

1. 每有一个新元素到来,布隆过滤器中的n递增

2. 调用 LocalBloomService.getKEncrypted方法,获取到元素val对应的k个bit位的偏移量offset

3. 通过offset>>5获取到bit位在[]int中的索引,思路同Exist

4. 通过 offset &31获取到bit位在int中的bit位置,思路同Exist

5. 通过|操作,将对应的bit位置为1

6. 重复上述流程,将k个bit位均置为1
*/
func (bloom *LocalBloomService) Set(val string) {
	if bloom.hasher == nil {
		bloom.hasher = Murmur3
	}
	bloom.N++
	for _, offset := range bloom.getKEncrypted(val) {
		index := offset >> 5
		bitOffset := offset & 31
		bloom.Bitmap[index] |= (1 << bitOffset)
	}
}
