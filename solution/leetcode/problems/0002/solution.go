package main

import (
    . "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
)

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) (ans *ListNode) {
    cur, sum := ans, 0
    for l1 != nil || l2 != nil || sum != 0 {
        if l1 != nil {
            sum += l1.Val
            l1 = l1.Next
        }
        if l2 != nil {
            sum += l2.Val
            l2 = l2.Next
        }
        if ans == nil {
            ans = &ListNode{Val: sum % 10}
            cur = ans
        } else {
            cur.Next = &ListNode{Val: sum % 10}
            cur = cur.Next
        }
        sum = sum / 10
    }
    return
}

// https://blog.lichangao.com/daily_practice/leetcode/list/array.html
