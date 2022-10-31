package functions

func gcd(a, b int) int {
    for a != 0 {
        a, b = b%a, a
    }
    return b
}
