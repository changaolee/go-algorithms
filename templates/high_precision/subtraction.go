package high_precision

// 高精度减法
// C = A - B，满足 A >= B, A >= 0, B >= 0
func sub(A, B []int) (C []int) {
    t := 0
    for i := 0; i < len(A); i++ {
        t = A[i] - t
        if i < len(B) {
            t -= B[i]
        }
        C = append(C, (t+10)%10)
        if t < 0 {
            t = 1
        } else {
            t = 0
        }
    }
    // 去除前导 0
    for len(C) > 1 && C[len(C)-1] == 0 {
        C = C[:len(C)-1]
    }
    return
}

// 比较 A 和 B 的大小关系，A >= B 时返回 true
func cmp(A, B []int) bool {
    if len(A) != len(B) {
        return len(A) > len(B)
    }
    for i := len(A) - 1; i >= 0; i-- {
        if A[i] != B[i] {
            return A[i] > B[i]
        }
    }
    return true
}
