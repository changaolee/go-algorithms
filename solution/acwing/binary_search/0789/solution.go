package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n, m int
    Fscan(in, &n, &m)

    q := make([]int, n)
    for i := 0; i < n; i++ {
        Fscan(in, &q[i])
    }
    var x int
    for i := 0; i < m; i++ {
        Fscan(in, &x)
        l, r := 0, n-1
        for l < r {
            mid := l + (r-l)/2
            if q[mid] >= x {
                r = mid
            } else {
                l = mid + 1
            }
        }
        if q[r] != x {
            Fprintf(out, "-1 -1\n")
        } else {
            Fprintf(out, "%d ", r)

            l, r = 0, n-1
            for l < r {
                mid := l + (r-l)/2 + 1
                if q[mid] <= x {
                    l = mid
                } else {
                    r = mid - 1
                }
            }
            Fprintf(out, "%d\n", l)
        }
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
