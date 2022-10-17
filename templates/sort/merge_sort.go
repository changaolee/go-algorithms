package sort

// 归并排序
func mergeSort(q []int, l int, r int) {
    if l >= r {
        return
    }
    mid := l + (r-l)/2
    mergeSort(q, l, mid)
    mergeSort(q, mid+1, r)

    tmp := make([]int, r-l+1)
    k, i, j := 0, l, mid+1
    for i <= mid && j <= r {
        if q[i] < q[j] {
            tmp[k] = q[i]
            k++
            i++
        } else {
            tmp[k] = q[j]
            k++
            j++
        }
    }
    for i <= mid {
        tmp[k] = q[i]
        k++
        i++
    }
    for j <= r {
        tmp[k] = q[j]
        k++
        j++
    }
    for i, idx := 0, l; idx <= r; i, idx = i+1, idx+1 {
        q[idx] = tmp[i]
    }
}
