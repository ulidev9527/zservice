package zservice

import (
	"math"
)

// 三维向量
type Vector3 struct {
	X float32 `json:"x,omitempty"`
	Y float32 `json:"y,omitempty"`
	Z float32 `json:"z,omitempty"`
}

// 计算到目标点的距离
func (v *Vector3) Distance(b Vector3) float32 {
	return float32(math.Sqrt(math.Pow(float64(b.X-v.X), 2) + math.Pow(float64(b.Y-v.Y), 2) + math.Pow(float64(b.Z-v.Z), 2)))
}
