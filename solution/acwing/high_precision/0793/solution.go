package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

// 高精度乘法
// // C = A * b，满足 A >= 0, b >= 0
func mul(A []int, b int) (C []int) {
    t := 0
    for i := 0; i < len(A) || t > 0; i++ {
        if i < len(A) {
            t += A[i] * b
        }
        C = append(C, t%10)
        t /= 10
    }
    // 去除前导 0
    for len(C) > 1 && C[len(C)-1] == 0 {
        C = C[:len(C)-1]
    }
    return
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var a string
    var b int
    Fscan(in, &a, &b)

    var A []int
    for i := len(a) - 1; i >= 0; i-- {
        A = append(A, int(a[i]-'0'))
    }

    C := mul(A, b)
    for i := len(C) - 1; i >= 0; i-- {
        Fprintf(out, "%d", C[i])
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
