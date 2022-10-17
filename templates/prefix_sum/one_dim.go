package prefix_sum

type PrefixSum1D struct {
    sum []int
}

// NewPrefixSum1D 初始化一维前缀和
func NewPrefixSum1D(q []int) *PrefixSum1D {
    n := len(q)
    sum := make([]int, n+1)
    for i := 1; i <= n; i++ {
        sum[i] = sum[i-1] + q[i-1]
    }
    return &PrefixSum1D{sum: sum}
}

// 查询一维前缀和
func (ps *PrefixSum1D) query(l, r int) int {
    return ps.sum[r] - ps.sum[l-1]
}
