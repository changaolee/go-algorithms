package cmd

import (
    "flag"
    "fmt"
    "github.com/changaolee/go-algorithms/pkg/leetcode"
    "os"
    "strconv"
    "testing"
)

func TestGenLCProblemCode(t *testing.T) {
    problemId, err := strconv.Atoi(flag.Arg(0))
    if err != nil {
        t.Fatal(err)
    }

    username := os.Getenv("LEETCODE_USERNAME_ZH")
    password := os.Getenv("LEETCODE_PASSWORD_ZH")

    dir := fmt.Sprintf("../solution/leetcode/problems/%04d/", problemId) // 自定义生成目录

    err = leetcode.GenLeetCodeSingleTest(username, password, problemId, false, dir, "")
    if err != nil {
        t.Fatal(err)
    }
}

func TestGenLCContestCode(t *testing.T) {
    args := flag.Args()
    id, err := strconv.Atoi(args[1])
    if err != nil {
        t.Fatal(err)
    }

    username := os.Getenv("LEETCODE_USERNAME_ZH")
    password := os.Getenv("LEETCODE_PASSWORD_ZH")

    var tag, dir string
    switch args[0] {
    case "weekly":
        contestId := leetcode.GetWeeklyContestId(id) // 自动生成下一场周赛 ID
        tag = leetcode.GetWeeklyContestTag(contestId)
        dir = fmt.Sprintf("../solution/leetcode/weekly/%d/", contestId) // 自定义生成目录
    case "biweekly":
        contestId := leetcode.GetBiweeklyContestId(id) // 自动生成下一场双周赛 ID
        tag = leetcode.GetBiweeklyContestTag(contestId)
        dir = fmt.Sprintf("../solution/leetcode/biweekly/%d/", contestId) // 自定义生成目录
    default:
        t.Fatal("error contest_type")
    }

    err = leetcode.GenLeetCodeTests(username, password, tag, false, dir, "")
    if err != nil {
        t.Fatal(err)
    }
}
