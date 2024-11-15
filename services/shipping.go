package services

import (
    "math"
)

func CalculateShippingCost(weight float64) int {
    chargeableWeight := math.Ceil(weight)
    if chargeableWeight > 100 {
        chargeableWeight = 100
    }

    if chargeableWeight <= 1 {
        return 18
    }

    cost := 18
    for i := 2.0; i <= chargeableWeight; i++ {
        cost += 5
        cost = int(math.Round(float64(cost) * 1.01))
    }

    return cost
}