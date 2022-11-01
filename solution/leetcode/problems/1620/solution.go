package main

import "math"

// 从左下角开始暴力枚举每个坐标，维护信号强度最大的坐标作为结果
func bestCoordinate(towers [][]int, radius int) (ans []int) {
    val := 0
    for i := 0; i <= 100; i++ {
        for j := 0; j <= 100; j++ {
            cur := 0
            for _, tower := range towers {
                x, y, q := tower[0], tower[1], tower[2]
                d := (x-i)*(x-i) + (y-j)*(y-j)
                if d <= radius*radius {
                    cur += int(float64(q) / (1 + math.Sqrt(float64(d))))
                }
            }
            if cur > val {
                val = cur
                ans = []int{i, j}
            }
        }
    }
    return
}
