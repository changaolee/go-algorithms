package main

import (
    "bufio"
    . "fmt"
    "io"
    "math"
    "os"
)

const eps = 1e-8

func run(_r io.Reader, out io.Writer) {
    in := bufio.NewReader(_r)

    var x float64
    Fscan(in, &x)

    minus := x < 0
    x = math.Abs(x)

    l, r := float64(0), math.Max(1, x)
    for r-l > eps {
        mid := l + (r-l)/2
        if mid*mid*mid <= x {
            l = mid
        } else {
            r = mid
        }
    }
    if minus {
        l = -l
    }
    Fprintf(out, "%.6f", l)
}

func main() {
    run(os.Stdin, os.Stdout)
}
