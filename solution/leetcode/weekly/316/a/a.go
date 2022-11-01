package main

// 直接比较字符串，没有交集的对立事件
func haveConflict(event1 []string, event2 []string) (ans bool) {
    ans = event1[1] >= event2[0] && event1[0] <= event2[1]
    return
}
