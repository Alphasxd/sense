package utils

import (
    "math"
    "math/rand"
)

func GenerateWeight() float64 {
    // 生成一个 0~1 之间的随机数并通过倒数的方式生成重量
    r := rand.Float64()
    weight := 1 / math.Pow(r, 2)
    // 限制最大重量为 100KG
    if weight > 100 {
        weight = 100
    }
    // 保留两位小数
    return math.Ceil(weight*100) / 100
}