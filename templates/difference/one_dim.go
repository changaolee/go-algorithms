package difference

type Difference1D struct {
    diff []int
    n    int
}

// NewDifference1D 初始化一维差分数组
func NewDifference1D(n int) *Difference1D {
    return &Difference1D{
        diff: make([]int, n+2),
        n:    n,
    }
}

// 给 l 到 r 区间内的元素加上 c
func (df *Difference1D) insert(l, r, c int) {
    df.diff[l] += c
    df.diff[r+1] -= c
}

// 获取最终数组
func (df *Difference1D) final() (ans []int) {
    ans = make([]int, df.n+1)
    for i := range ans {
        ans[i] = df.diff[i]
    }
    for i := 1; i <= df.n; i++ {
        ans[i] += ans[i-1]
    }
    ans = ans[1:]
    return
}
