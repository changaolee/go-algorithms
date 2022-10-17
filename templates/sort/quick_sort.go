package sort

// 快速排序
func quickSort(q []int, l int, r int) {
    if l >= r {
        return
    }
    x := q[l+(r-l)/2]
    i, j := l-1, r+1
    for i < j {
        for {
            i++
            if q[i] >= x {
                break
            }
        }
        for {
            j--
            if q[j] <= x {
                break
            }
        }
        if i < j {
            q[i], q[j] = q[j], q[i]
        }
    }
    quickSort(q, l, j)
    quickSort(q, j+1, r)
}
