package main

func lengthOfLongestSubstring(s string) (ans int) {
    pos := map[byte]int{}
    for i, j := 0, 0; i < len(s); i++ {
        if p, ok := pos[s[i]]; ok {
            j = max(j, p+1)
        }
        ans = max(ans, i-j+1)
        pos[s[i]] = i
    }
    return
}

func max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

// https://blog.lichangao.com/daily_practice/leetcode/string/sliding_window.html
