package zservice

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"github.com/rs/xid"
)

const charsetLow = "abcdefghijklmnopqrstuvwxyz"
const charsetUp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetAll = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机数字 int
func RandomInt(max int) int {
	return RandomIntRange(0, max)
}

// 随机数字范围 int
func RandomIntRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// 随机数字 int
func RandomInt32(max int32) int32 {
	return int32(RandomInt(int(max)))
}

// 随机数字范围 int
func RandomInt32Range(min int32, max int32) int32 {
	return int32(RandomIntRange(int(min), int(max)))
}

// 随机数字 uint32
func RandomUInt32(max uint32) uint32 {
	return RandomUInt32Range(0, max)
}

// 随机数字范围 uint32
func RandomUInt32Range(min uint32, max uint32) uint32 {
	return rand.Uint32()%(max-min+1) + min
}

// 随机数字
func RandomInt64(max int64) int64 {
	return RandomInt64Range(0, max)
}

// 随机数字范围 int64
func RandomInt64Range(min int64, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomInt64RangeAS(min int64, max int64) int64 {
	if RandomInt(1) == 0 {
		return RandomInt64Range(min, max)
	} else {
		return -RandomInt64Range(min, max)
	}
}

func RandomFloat32(max float32) float32 {
	return RandomFloat32Range(0, max)
}

// 随机数字范围
func RandomFloat32Range(min float32, max float32) float32 {
	return rand.Float32()*(max-min) + min
}

// 随机 float 正负
func RandomFloat32RangeAS(min float32, max float32) float32 {
	if RandomInt(1) == 0 {
		return RandomFloat32Range(min, max)
	} else {
		return -RandomFloat32Range(min, max)
	}
}

func RandomFloat64RangeAS(min float64, max float64) float64 {
	if RandomInt(1) == 0 {
		return RandomFloat64Range(min, max)
	} else {
		return -RandomFloat64Range(min, max)
	}
}

// 随机数字范围
func RandomFloat64Range(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// 随机概率 - 返回随机到的索引
func RandomRateInt(arr []int) int {
	total := 0
	for _, val := range arr {
		total += val
	}
	random := rand.Intn(total)
	for i, val := range arr {
		if random < val {
			return i
		}
		random -= val
	}
	return len(arr) - 1 // 兜底逻辑
}

func RandomRateFloat32(arr []float32) int {
	total := 0.0
	for _, val := range arr {
		total += float64(val)
	}
	random := rand.Float64() * total
	for i, val := range arr {
		if random < float64(val) {
			return i
		}
		random -= float64(val)
	}
	return len(arr) - 1 // 兜底逻辑
}

// 含浮点数的随机概率 - 返回随机到的索引
func RandomRateFloat64(arr []float64) int {
	total := 0.0
	for _, val := range arr {
		total += val
	}
	random := rand.Float64() * total
	for i, val := range arr {
		if random < val {
			return i
		}
		random -= val
	}
	return len(arr) - 1 // 兜底逻辑
}

// 含浮点数的随机概率 - 返回随机到的索引
func RandomRateInt64(arr []int64) int {
	fArr := []float64{}
	for _, v := range arr {
		fArr = append(fArr, float64(v))
	}
	return RandomRateFloat64(fArr)
}

// 随机字符串
// @count 字符串长度
// @charset 字符集
func randomString_base(count int, charset string) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charset[RandomInt(len(charset)-1)]
	}
	return string(b)
}

// 随机字符串
// @count 字符串长度
func RandomString(count int) string {
	return randomString_base(count, charsetAll)
}

// 随机大写字符串
func RandomStringUP(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charsetUp[RandomInt(len(charsetUp)-1)]
	}
	return string(b)
}

// 随机小写字符串
func RandomStringLow(count int) string {
	b := make([]byte, count)
	for i := range b {
		b[i] = charsetLow[RandomInt(len(charsetLow)-1)]
	}
	return string(b)
}

// md5
func RandomMD5() string {
	return RandomMD5_XID_Random()
}

// md5 xid
func RandomMD5_XID() string {
	return hex.EncodeToString(xid.New().Bytes())
}

// md5 xid + random string
func RandomMD5_XID_Random() string {

	m := md5.New()
	m.Write(append(xid.New().Bytes(), []byte(RandomString(32))...))
	return hex.EncodeToString(m.Sum(nil))
}

// 随机 xid
func RandomXID() string {
	return xid.New().String()
}

// 随机bool
func RandomBool() bool { return RandomInt(9)%2 == 1 }
