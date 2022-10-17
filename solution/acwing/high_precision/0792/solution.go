package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

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

    var C []int
    if cmp(A, B) {
        C = sub(A, B)
    } else {
        C = sub(B, A)
        Fprintf(out, "-")
    }
    for i := len(C) - 1; i >= 0; i-- {
        Fprintf(out, "%d", C[i])
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
