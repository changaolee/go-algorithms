package main

func countSubarrays(a []int, minK int, maxK int) int64 {
    ans := 0
    minI, maxI, i0 := -1, -1, -1
    for i, x := range a {
        if x == minK {
            minI = i
        }
        if x == maxK {
            maxI = i
        }
        if x < minK || x > maxK {
            i0 = i
        }
        if min(minI, maxI) > i0 {
            ans += min(minI, maxI) - i0
        }
    }
    return int64(ans)
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}
