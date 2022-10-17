package main

import "math"

func reverse(x int) (ans int) {
    for x != 0 {
        if ans > 0 && ans > (math.MaxInt32-x%10)/10 {
            ans = 0
            break
        }
        if ans < 0 && ans < (math.MinInt32-x%10)/10 {
            ans = 0
            break
        }
        ans = ans*10 + x%10
        x /= 10
    }
    return
}
