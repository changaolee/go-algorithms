package main

import (
    "testing"

    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
)

func TestB(t *testing.T) {
    targetCaseNum := 0 // -1
    if err := testutils.RunLeetCodeFuncWithFile(t, subarrayGCD, "b.txt", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// https://leetcode.cn/contest/weekly-contest-316/problems/number-of-subarrays-with-gcd-equal-to-k/
