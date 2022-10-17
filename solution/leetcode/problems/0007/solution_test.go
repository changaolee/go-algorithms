package main

import (
    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
    "testing"
)

func Test0007(t *testing.T) {
    targetCaseNum := 0 // -1
    if err := testutils.RunLeetCodeFuncWithFile(t, reverse, "sample.txt", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// https://leetcode.cn/problems/reverse-integer/
