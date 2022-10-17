package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

type PrefixSum1D struct {
    sum []int
}

// NewPrefixSum1D 初始化一维前缀和
func NewPrefixSum1D(q []int) *PrefixSum1D {
    n := len(q)
    sum := make([]int, n+1)
    for i := 1; i <= n; i++ {
        sum[i] = sum[i-1] + q[i-1]
    }
    return &PrefixSum1D{sum: sum}
}

// 查询一维前缀和
func (ps *PrefixSum1D) query(l, r int) int {
    return ps.sum[r] - ps.sum[l-1]
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n, m int
    Fscan(in, &n, &m)

    q := make([]int, n)
    for i := 0; i < n; i++ {
        Fscan(in, &q[i])
    }

    ps := NewPrefixSum1D(q)

    for i := 0; i < m; i++ {
        var l, r int
        Fscan(in, &l, &r)
        Fprintf(out, "%d\n", ps.query(l, r))
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
