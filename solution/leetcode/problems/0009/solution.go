package main

import "math"

func isPalindrome(x int) bool {
    if x < 0 {
        return false
    }
    num, cur := x, 0
    for x != 0 {
        if cur > (math.MaxInt32-x%10)/10 {
            return false
        }
        cur = cur*10 + x%10
        x /= 10
    }
    return num == cur
}
