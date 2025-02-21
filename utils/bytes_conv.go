package utils

import "unsafe"

// Bytes2String 将字节切片无拷贝地转化为字符串
func Bytes2String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
