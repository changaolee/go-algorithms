package main

import (
    "testing"

    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
)

func Test0009(t *testing.T) {
    targetCaseNum := 0 // -1
    if err := testutils.RunLeetCodeFuncWithFile(t, isPalindrome, "sample.txt", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// https://leetcode.cn/problems/palindrome-number/
