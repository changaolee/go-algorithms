package leetcode

import (
    "strconv"
    "time"
)

// id > 0，指定具体的一场周赛
// id = 0，指定下一场或当前正在进行的周赛
// id < 0，指定上 |id| 场周赛（例如 id = -1 表示最近的一场结束的周赛）

func GetWeeklyContestId(contestId int) int {
    if contestId <= 0 {
        utc8, err := time.LoadLocation("Asia/Shanghai")
        if err != nil {
            panic(err)
        }

        // 以 2020 年第一场周赛的结束时间为基准
        endTime170 := time.Date(2020, 1, 5, 12, 0, 0, 0, utc8)
        weeksSince170 := 1 + int(time.Since(endTime170)/(7*24*time.Hour))
        contestId += 170 + weeksSince170
    }
    return contestId
}

func GetBiweeklyContestId(contestId int) int {
    if contestId <= 0 {
        utc8, err := time.LoadLocation("Asia/Shanghai")
        if err != nil {
            panic(err)
        }

        // 以 2020 年第一场双周赛的结束时间为基准
        endTime17 := time.Date(2020, 1, 12, 0, 0, 0, 0, utc8)
        twoWeeksSince17 := 1 + int(time.Since(endTime17)/(14*24*time.Hour))
        contestId += 17 + twoWeeksSince17
    }
    return contestId
}

func GetWeeklyContestTag(contestId int) string {
    return "weekly-contest-" + strconv.Itoa(GetWeeklyContestId(contestId))
}

func GetBiweeklyContestTag(contestId int) string {
    return "biweekly-contest-" + strconv.Itoa(GetBiweeklyContestId(contestId))
}
