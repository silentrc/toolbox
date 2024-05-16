package utils

import (
	"math/rand"
	"time"
)

type generateUtils struct {
}

// 生成类
func (u *utils) NewGenerateUtils() *generateUtils {
	return &generateUtils{}
}

// 生成随机字符串
func (g *generateUtils) RandString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	newBytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, newBytes[r.Intn(len(newBytes))])
	}
	return string(result)
}
