package main

func longestPalindrome(s string) (ans string) {
    for i := 0; i < len(s); i++ {
        // 奇数长度
        l, r := i-1, i+1
        for l >= 0 && r < len(s) && s[l] == s[r] {
            l--
            r++
        }
        if len(ans) < r-l-1 {
            ans = s[l+1 : r]
        }

        // 偶数长度
        l, r = i, i+1
        for l >= 0 && r < len(s) && s[l] == s[r] {
            l--
            r++
        }
        if len(ans) < r-l-1 {
            ans = s[l+1 : r]
        }
    }
    return
}
