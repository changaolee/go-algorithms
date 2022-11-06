package main

import "strings"

func interpret(command string) (ans string) {
    ans = strings.ReplaceAll(command, "()", "o")
    ans = strings.ReplaceAll(ans, "(al)", "al")
    return
}
