package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

type Difference1D struct {
    diff []int
    n    int
}

// NewDifference1D 初始化一维差分数组
func NewDifference1D(n int) *Difference1D {
    return &Difference1D{
        diff: make([]int, n+2),
        n:    n,
    }
}

// 给 l 到 r 区间内的元素加上 c
func (df *Difference1D) insert(l, r, c int) {
    df.diff[l] += c
    df.diff[r+1] -= c
}

// 获取最终数组
func (df *Difference1D) final() (ans []int) {
    ans = make([]int, df.n+1)
    for i := range ans {
        ans[i] = df.diff[i]
    }
    for i := 1; i <= df.n; i++ {
        ans[i] += ans[i-1]
    }
    ans = ans[1:]
    return
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n, m int
    Fscan(in, &n, &m)

    df := NewDifference1D(n)

    var x int
    for i := 0; i < n; i++ {
        Fscan(in, &x)
        df.insert(i+1, i+1, x)
    }

    for i := 0; i < m; i++ {
        var l, r, c int
        Fscan(in, &l, &r, &c)
        df.insert(l, r, c)
    }

    ans := df.final()

    for _, num := range ans {
        Fprintf(out, "%d ", num)
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
