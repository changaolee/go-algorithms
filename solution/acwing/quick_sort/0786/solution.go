package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

func quickSelect(q []int, l int, r int, k int) int {
    if l == r {
        return q[l]
    }
    x := q[l+(r-l)/2]
    i, j := l-1, r+1
    for i < j {
        for {
            i++
            if q[i] >= x {
                break
            }
        }
        for {
            j--
            if q[j] <= x {
                break
            }
        }
        if i < j {
            q[i], q[j] = q[j], q[i]
        }
    }
    // 根据左半边区间元素数量，判断第 K 个数在哪个区间
    sl := j - l + 1
    if sl >= k {
        return quickSelect(q, l, j, k)
    } else {
        return quickSelect(q, j+1, r, k-sl)
    }
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n, k int
    Fscan(in, &n, &k)

    q := make([]int, n)
    for i := 0; i < n; i++ {
        Fscan(in, &q[i])
    }

    Fprintf(out, "%d", quickSelect(q, 0, n-1, k))
}

func main() {
    run(os.Stdin, os.Stdout)
}
