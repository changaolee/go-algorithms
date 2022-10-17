package main

func twoSum(a []int, trg int) (ans []int) {
    h := map[int]int{}
    for i, x := range a {
        if p, ok := h[trg-x]; ok {
            return []int{p, i}
        }
        h[x] = i
    }
    return
}

// https://blog.lichangao.com/daily_practice/leetcode/array/sum.html
