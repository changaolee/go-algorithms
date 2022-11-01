package main

// 暴力枚举
// 枚举以下标 i 为左边界的子数组，在不包含因子 k 时提前退出循环。
func subarrayGCD(a []int, k int) (ans int) {
    for i := range a {
        g := 0
        for _, x := range a[i:] {
            g = gcd(g, x)
            if g%k > 0 {
                break
            }
            if g == k {
                ans++
            }
        }
    }
    return
}

func gcd(a, b int) int {
    for a != 0 {
        a, b = b%a, a
    }
    return b
}
