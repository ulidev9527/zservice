package zservice

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/rs/xid"
)

const charsetLow = "abcdefghijklmnopqrstuvwxyz"
const charsetUp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetAll = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 随机数字
func RandomInt(count int) int {
	return seededRand.Intn(count)
}

// 随机数字
func RandomInt64(count int64) int64 {
	return seededRand.Int63n(count)
}

// 随机数字范围
func RandomIntRange(min int, max int) int {
	return RandomInt(max-min) + min
}

// 随机数字范围
func RandomInt64Range(min int64, max int64) int64 {
	return RandomInt64(max-min) + min
}

// 随机数字范围
func RandomFloat32Range(min float32, max float32) float32 {
	return seededRand.Float32()*(max-min) + min
}

// 随机数字范围
func RandomFloat64Range(min float64, max float64) float64 {
	return seededRand.Float64()*(max-min) + min
}

// 随机字符串
func RandomString(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charsetAll[RandomInt(len(charsetAll))]
	}
	return string(b)
}

// 随机大写字符串
func RandomStringUP(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charsetUp[RandomInt(len(charsetUp))]
	}
	return string(b)
}

// 随机小写字符串
func RandomStringLow(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charsetLow[RandomInt(len(charsetLow))]
	}
	return string(b)
}

// md5
func RandomMD5() string {
	m := md5.New()
	m.Write(xid.New().Bytes())
	return hex.EncodeToString(m.Sum(nil))
}

// 随机 xid
func RandomXID() string {
	return xid.New().String()
}
