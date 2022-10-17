package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

type PrefixSum2D struct {
    sum [][]int
}

// NewPrefixSum2D 初始化二维前缀和
func NewPrefixSum2D(q [][]int) *PrefixSum2D {
    n, m := len(q), len(q[0])
    sum := make([][]int, n+1)
    for i := range sum {
        sum[i] = make([]int, m+1)
    }
    for i := 1; i <= n; i++ {
        for j := 1; j <= m; j++ {
            sum[i][j] = sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1] + q[i-1][j-1]
        }
    }
    return &PrefixSum2D{sum: sum}
}

// 查询二维前缀和
func (ps *PrefixSum2D) query(x1, y1, x2, y2 int) int {
    return ps.sum[x2][y2] - ps.sum[x2][y1-1] - ps.sum[x1-1][y2] + ps.sum[x1-1][y1-1]
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n, m, k int
    Fscan(in, &n, &m, &k)

    q := make([][]int, n)
    for i := range q {
        q[i] = make([]int, m)
    }
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            Fscan(in, &q[i][j])
        }
    }

    ps := NewPrefixSum2D(q)

    for i := 0; i < k; i++ {
        var x1, y1, x2, y2 int
        Fscan(in, &x1, &y1, &x2, &y2)
        Fprintf(out, "%d\n", ps.query(x1, y1, x2, y2))
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
