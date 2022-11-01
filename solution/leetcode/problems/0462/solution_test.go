package main

import (
    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
    "testing"
)

func Test0462(t *testing.T) {
    targetCaseNum := 0 // -1
    if err := testutils.RunLeetCodeFuncWithFile(t, minMoves2, "sample.txt", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// https://leetcode.cn/problems/minimum-moves-to-equal-array-elements-ii/
