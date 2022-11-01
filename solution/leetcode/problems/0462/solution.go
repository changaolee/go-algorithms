package main

import "sort"

// 排序后，把中位数作为目标值，累加变化量。
// 注：可通过绝对值不等式证明该结论。
func minMoves2(a []int) (ans int) {
    sort.Ints(a)
    mid := a[len(a)/2]
    for _, x := range a {
        ans += abs(x - mid)
    }
    return
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
