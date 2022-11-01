package functions

func max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
