package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

func quickSort(q []int, l int, r int) {
    if l >= r {
        return
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
    quickSort(q, l, j)
    quickSort(q, j+1, r)
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n int
    Fscan(in, &n)

    q := make([]int, n)
    for i := 0; i < n; i++ {
        Fscan(in, &q[i])
    }

    quickSort(q, 0, n-1)

    for i := 0; i < n; i++ {
        Fprintf(out, "%d ", q[i])
    }
}

func main() {
    run(os.Stdin, os.Stdout)
}
