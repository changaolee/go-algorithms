package search

// 整数二分：模板 1
// 当 l = mid 时，mid 整除后 +1
func binarySearchInt1(l, r int) int {
    for l < r {
        mid := l + (r-l)/2 + 1
        if checkInt(mid) {
            l = mid
        } else {
            r = mid - 1
        }
    }
    return l
}

// 整数二分：模板 2
// 当 r = mid 时，mid 直接整除即可
func binarySearchInt2(l, r int) int {
    for l < r {
        mid := l + (r-l)/2
        if checkInt(mid) {
            r = mid
        } else {
            l = mid + 1
        }
    }
    return r
}

func checkInt(mid int) bool {
    // 检查逻辑
    return true
}

// 浮点数二分
func binarySearchFloat(l, r float64) float64 {
    const eps = 1e-8
    for r-l > eps {
        mid := l + (r-l)/2
        if checkFloat(mid) {
            l = mid
        } else {
            r = mid
        }
    }
    return l
}

func checkFloat(mid float64) bool {
    // 检查逻辑
    return true
}
