package main

func countDistinctIntegers(a []int) (ans int) {
    has := map[int]bool{}
    for _, x := range a {
        has[x] = true
        has[reverseNum(x)] = true
    }
    ans = len(has)
    return
}

func reverseNum(x int) (ans int) {
    for x > 0 {
        ans = ans*10 + x%10
        x /= 10
    }
    return
}
