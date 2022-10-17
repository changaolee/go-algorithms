package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

// 高精度除法
// A / b = C ... r，满足 A >= 0, b > 0
func div(A []int, b int) (C []int, r int) {
    for i := len(A) - 1; i >= 0; i-- {
        r = r*10 + A[i]
        C = append(C, r/b)
        r %= b
    }
    // reverse
    for i, j := 0, len(C)-1; i < j; i, j = i+1, j-1 {
        C[i], C[j] = C[j], C[i]
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
    C, r := div(A, b)
    for i := len(C) - 1; i >= 0; i-- {
        Fprintf(out, "%d", C[i])
    }
    Fprintf(out, "\n%d", r)
}

func main() {
    run(os.Stdin, os.Stdout)
}
