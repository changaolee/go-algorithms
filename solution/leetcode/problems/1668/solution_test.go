package main

import (
    "testing"

    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
)

func Test1668(t *testing.T) {
    targetCaseNum := 0 // -1
    if err := testutils.RunLeetCodeFuncWithFile(t, maxRepeating, "sample.txt", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// https://leetcode.cn/problems/maximum-repeating-substring/
