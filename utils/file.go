package utils

import "os"

// PathExists 判断文件是否存在
//
// 存在返回true, 否则返回false
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
