package high_precision

// 高精度加法
// C = A + B，满足 A >= 0, B >= 0
func add(A, B []int) (C []int) {
    t := 0
    for i := 0; i < len(A) || i < len(B); i++ {
        if i < len(A) {
            t += A[i]
        }
        if i < len(B) {
            t += B[i]
        }
        C = append(C, t%10)
        t /= 10
    }
    if t > 0 {
        C = append(C, t)
    }
    return
}
