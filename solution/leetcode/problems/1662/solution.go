package main

// 方法一：拼接成字符串后比较
// 时间复杂度：O(m + n)
// 空间复杂度：O(m + n)
//func arrayStringsAreEqual(x []string, y []string) bool {
//    return strings.Join(x, "") == strings.Join(y, "")
//}

// 方法二：逐个字符比较比较
// 时间复杂度：O(m + n)
// 空间复杂度：O(1)
func arrayStringsAreEqual(x []string, y []string) bool {
    var p, q, i, j int
    for p < len(x) && q < len(y) {
        if x[p][i] != y[q][j] {
            return false
        }
        i++
        if i == len(x[p]) {
            p++
            i = 0
        }
        j++
        if j == len(y[q]) {
            q++
            j = 0
        }
    }
    return p == len(x) && q == len(y)
}
