package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

// 高精度加法
// C = A + B, A >= 0, B >= 0
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

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var a, b string
    Fscan(in, &a, &b)

    var A, B []int
    for i := len(a) - 1; i >= 0; i-- {
        A = append(A, int(a[i]-'0'))
    }
    for i := len(b) - 1; i >= 0; i-- {
        B = append(B, int(b[i]-'0'))
    }

    C := add(A, B)
    for i := len(C) - 1; i >= 0; i-- {
        Fprintf(out, "%d", C[i])
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
