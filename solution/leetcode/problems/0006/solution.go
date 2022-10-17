package main

func convert(s string, numRows int) (ans string) {
    if numRows == 1 {
        ans = s
        return
    }
    rows := make([]string, numRows)
    row, down := 0, false
    for _, c := range s {
        rows[row] += string(c)
        if row == 0 || row == numRows-1 {
            down = !down
        }
        if down {
            row++
        } else {
            row--
        }
    }
    for i := range rows {
        ans += rows[i]
    }
    return
}
