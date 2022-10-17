package difference

type Difference2D struct {
    diff [][]int
    n, m int
}

// NewDifference2D 初始化二维差分数组
func NewDifference2D(n, m int) *Difference2D {
    diff := make([][]int, n+2)
    for i := range diff {
        diff[i] = make([]int, m+2)
    }
    return &Difference2D{
        diff: diff,
        n:    n,
        m:    m,
    }
}

// 给子矩阵内的元素加上 c
// 子矩阵的左上角和右下角坐标为 (x1, y1) 和 (x2, y2)
func (df *Difference2D) insert(x1, y1, x2, y2, c int) {
    df.diff[x1][y1] += c
    df.diff[x2+1][y1] -= c
    df.diff[x1][y2+1] -= c
    df.diff[x2+1][y2+1] += c
}

// 获取最终矩阵
func (df *Difference2D) final() (ans [][]int) {
    ans = make([][]int, df.n+1)
    for i := range ans {
        ans[i] = make([]int, df.m+1)
        for j := range ans[i] {
            ans[i][j] = df.diff[i][j]
        }
    }
    for i := 1; i <= df.n; i++ {
        for j := 1; j <= df.m; j++ {
            ans[i][j] += ans[i-1][j] + ans[i][j-1] - ans[i-1][j-1]
        }
    }
    for i := range ans {
        ans[i] = ans[i][1:]
    }
    ans = ans[1:]
    return
}
