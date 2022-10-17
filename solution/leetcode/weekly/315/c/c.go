package main

func sumOfNumberAndReverse(n int) (ans bool) {
    for i := 0; i <= n; i++ {
        if i+reverseNum(i) == n {
            ans = true
            return
        }
    }
    return
}

func reverseNum(x int) (ans int) {
    for x > 0 {
        ans = ans*10 + x%10
        x /= 10
    }
    return
}
