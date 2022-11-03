package main

import "strings"

func maxRepeating(sequence string, s string) (ans int) {
    for k := len(sequence) / len(s); k > 0; k-- {
        if strings.Contains(sequence, strings.Repeat(s, k)) {
            return k
        }
    }
    return
}
