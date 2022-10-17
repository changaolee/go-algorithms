package main

func findMaxK(a []int) (ans int) {
    ans = -1
    has := map[int]bool{}
    for _, x := range a {
        if abs(x) > ans && has[-x] {
            ans = abs(x)
        }
        has[x] = true
    }
    return
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
