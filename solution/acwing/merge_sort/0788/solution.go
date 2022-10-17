package main

import (
    "bufio"
    . "fmt"
    "io"
    "os"
)

func mergeSort(q []int, l int, r int) (ans int) {
    if l >= r {
        return
    }
    mid := l + (r-l)/2
    ans += mergeSort(q, l, mid)
    ans += mergeSort(q, mid+1, r)

    tmp := make([]int, r-l+1)
    k, i, j := 0, l, mid+1
    for i <= mid && j <= r {
        if q[i] <= q[j] {
            tmp[k] = q[i]
            k++
            i++
        } else {
            ans += mid - i + 1
            tmp[k] = q[j]
            k++
            j++
        }
    }
    for i <= mid {
        tmp[k] = q[i]
        k++
        i++
    }
    for j <= r {
        tmp[k] = q[j]
        k++
        j++
    }
    for i, idx := 0, l; idx <= r; i, idx = i+1, idx+1 {
        q[idx] = tmp[i]
    }
    return
}

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var n int
    Fscan(in, &n)

    q := make([]int, n)
    for i := 0; i < n; i++ {
        Fscan(in, &q[i])
    }

    Fprintf(out, "%d", mergeSort(q, 0, n-1))
}

func main() {
    run(os.Stdin, os.Stdout)
}
