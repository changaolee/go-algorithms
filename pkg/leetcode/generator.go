package leetcode

import (
    "fmt"
    "github.com/levigross/grequests"
    "github.com/skratchdot/open-golang/open"
    "net/url"
    "strconv"
    "strings"
    "sync"
    "time"
)

const host = "leetcode.cn"
const graphqlURL = "https://" + host + "/graphql"

// 使用用户名和密码登录
func login(username, password string) (session *grequests.Session, err error) {
    const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"
    session = grequests.NewSession(&grequests.RequestOptions{
        UserAgent:    ua,
        UseCookieJar: true,
    })

    // "touch" csrfToken
    resp, err := session.Post(graphqlURL, &grequests.RequestOptions{
        JSON: map[string]interface{}{
            "operationName": "globalData",
            "query":         "query globalData {\n  feature {\n    questionTranslation\n    subscription\n    signUp\n    discuss\n    mockInterview\n    contest\n    store\n    book\n    chinaProblemDiscuss\n    socialProviders\n    studentFooter\n    cnJobs\n    __typename\n  }\n  userStatus {\n    isSignedIn\n    isAdmin\n    isStaff\n    isSuperuser\n    isTranslator\n    isPremium\n    isVerified\n    isPhoneVerified\n    isWechatVerified\n    checkedInToday\n    username\n    realName\n    userSlug\n    groups\n    jobsCompany {\n      nameSlug\n      logo\n      description\n      name\n      legalName\n      isVerified\n      permissions {\n        canInviteUsers\n        canInviteAllSite\n        leftInviteTimes\n        maxVisibleExploredUser\n        __typename\n      }\n      __typename\n    }\n    avatar\n    optedIn\n    requestRegion\n    region\n    activeSessionId\n    permissions\n    notificationStatus {\n      lastModified\n      numUnread\n      __typename\n    }\n    completedFeatureGuides\n    useTranslation\n    __typename\n  }\n  siteRegion\n  chinaHost\n  websocketUrl\n}\n",
        },
    })
    if err != nil {
        // maybe timeout
        fmt.Println("访问失败，重试", err)
        time.Sleep(time.Second)
        return login(username, password)
    }
    if !resp.Ok {
        return nil, fmt.Errorf("POST %s return code %d", graphqlURL, resp.StatusCode)
    }

    var csrfToken string
    for _, c := range resp.RawResponse.Cookies() {
        if c.Name == "csrftoken" {
            csrfToken = c.Value
            break
        }
    }
    if csrfToken == "" {
        return nil, fmt.Errorf("csrftoken not found in response")
    }

    // log in
    loginURL := fmt.Sprintf("https://%s/accounts/login/", host)
    resp, err = session.Post(loginURL, &grequests.RequestOptions{
        Data: map[string]string{
            "csrfmiddlewaretoken": csrfToken,
            "login":               username,
            "password":            password,
            "next":                "/",
        },
        Headers: map[string]string{
            "origin":  "https://" + host,
            "referer": "https://" + host + "/",
        },
    })
    if err != nil {
        fmt.Println("访问失败，重试", err)
        time.Sleep(time.Second)
        return login(username, password)
    }
    if !resp.Ok {
        return nil, fmt.Errorf("POST %s return code %d", loginURL, resp.StatusCode)
    }

    u, _ := url.Parse(loginURL)
    for _, cookie := range session.HTTPClient.Jar.Cookies(u) {
        if cookie.Name == "LEETCODE_SESSION" {
            return
        }
    }

    return nil, fmt.Errorf("登录失败：账号或密码错误")
}

func fetchSingleProblemURL(session *grequests.Session, problemId int) (*problem, error) {
    resp, err := session.Post(graphqlURL, &grequests.RequestOptions{
        JSON: map[string]any{
            "query": "\n    query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {\n  problemsetQuestionList(\n    categorySlug: $categorySlug\n    limit: $limit\n    skip: $skip\n    filters: $filters\n  ) {\n    hasMore\n    total\n    questions {\n      acRate\n      difficulty\n      freqBar\n      frontendQuestionId\n      isFavor\n      paidOnly\n      solutionNum\n      status\n      title\n      titleCn\n      titleSlug\n      topicTags {\n        name\n        nameTranslated\n        id\n        slug\n      }\n      extra {\n        hasVideoSolution\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n  }\n}\n    ",
            "variables": map[string]any{
                "categorySlug": "",
                "skip":         0,
                "limit":        50,
                "filters": map[string]any{
                    "searchKeywords": strconv.Itoa(problemId),
                },
            },
        },
    })
    if err != nil {
        fmt.Println("获取题目链接失败，重试", err)
        time.Sleep(time.Second)
        return fetchSingleProblemURL(session, problemId)
    }
    if !resp.Ok {
        return nil, fmt.Errorf("POST %s return code %d", graphqlURL, resp.StatusCode)
    }

    d := struct {
        Data struct {
            ProblemsetQuestionList struct {
                Questions []struct {
                    QuestionId string `json:"frontendQuestionId"`
                    TitleSlug  string `json:"titleSlug"`
                } `json:"questions"`
            } `json:"problemsetQuestionList"`
        } `json:"data"`
    }{}
    if err = resp.JSON(&d); err != nil {
        return nil, err
    }

    for _, q := range d.Data.ProblemsetQuestionList.Questions {
        if qid, err := strconv.Atoi(q.QuestionId); err == nil && qid == problemId {
            p := &problem{
                id:            q.QuestionId,
                url:           fmt.Sprintf("https://%s/problems/%s/", host, q.TitleSlug),
                isFuncProblem: true,
            }
            return p, nil
        }
    }
    return nil, fmt.Errorf("search problem %d error", problemId)
}

// 获取题目信息（含题目链接）
// contestTag 如 "weekly-contest-200"，可以从比赛链接中获取
func fetchProblemURLs(session *grequests.Session, contestTag string) (problems []*problem, err error) {
    contestInfoURL := fmt.Sprintf("https://%s/contest/api/info/%s/", host, contestTag)
    resp, err := session.Get(contestInfoURL, nil)
    if err != nil {
        return
    }
    if !resp.Ok {
        return nil, fmt.Errorf("GET %s return code %d", contestInfoURL, resp.StatusCode)
    }

    d := struct {
        Contest struct {
            ID              int    `json:"id"`
            OriginStartTime int64  `json:"origin_start_time"`
            StartTime       int64  `json:"start_time"`
            Title           string `json:"title"`
        } `json:"contest"`
        Questions []struct {
            Credit    int    `json:"credit"`     // 得分/难度
            Title     string `json:"title"`      // 题目标题
            TitleSlug string `json:"title_slug"` // 题目链接
        } `json:"questions"`
        Registered bool `json:"registered"` // 是否报名
        UserNum    int  `json:"user_num"`   // 参赛人数
    }{}
    if err = resp.JSON(&d); err != nil {
        return
    }
    if d.Contest.StartTime == 0 {
        return nil, fmt.Errorf("未找到比赛或比赛尚未开始: %s", contestTag)
    }

    //fmt.Println("当前报名人数", d.UserNum)

    if sleepTime := time.Until(time.Unix(d.Contest.StartTime, 0)); sleepTime > 0 {
        if !d.Registered {
            fmt.Printf("该账号尚未报名%s\n", d.Contest.Title)
            return
        }

        sleepTime += 500 * time.Millisecond // 消除误差
        fmt.Printf("%s尚未开始，等待中……\n%v\n", d.Contest.Title, sleepTime)
        time.Sleep(sleepTime)
        return fetchProblemURLs(session, contestTag)
    }

    if len(d.Questions) == 0 {
        return nil, fmt.Errorf("题目链接为空: %s", contestTag)
    }

    fmt.Println("难度 标题")
    for _, q := range d.Questions {
        fmt.Printf("%3d %s\n", q.Credit, q.Title)
    }

    problems = make([]*problem, len(d.Questions))
    for i, q := range d.Questions {
        problems[i] = &problem{
            id:  string(byte('a' + i)),
            url: fmt.Sprintf("https://%s/contest/%s/problems/%s/", host, contestTag, q.TitleSlug),

            isFuncProblem: true,
        }
    }
    return
}

func updateComment(cmt string) string {
    if cmt != "" {
        if !strings.HasPrefix(cmt, "//") {
            cmt = "// " + cmt
        }
        cmt += "\n"
    }
    return cmt
}

func handleContestProblems(session *grequests.Session, problems []*problem) error {
    wg := &sync.WaitGroup{}
    wg.Add(1 + len(problems))

    go func() {
        defer wg.Done()
        for _, p := range problems {
            if p.openURL {
                if err := open.Run(p.url); err != nil {
                    fmt.Println("open err:", p.url, err)
                }
            }
        }
    }()

    for _, p := range problems {
        fmt.Println(p.id, p.url)

        go func(p *problem) {
            defer wg.Done()

            err := parseContestProblem(session, p)
            if err != nil {
                return
            }
        }(p)
    }

    wg.Wait()
    return nil
}

func parseContestProblem(session *grequests.Session, p *problem) error {
    if err := p.parseHTML(session); err != nil {
        fmt.Println(err)
        return nil
    }

    customFuncContent := "\t\n\treturn"
    if p.needMod {
        customFuncContent = "\t\n\t\n\t\n\tans = (ans%mod + mod) % mod\n\treturn"
    }

    p.defaultCode = modifyDefaultCode(p.defaultCode, p.funcLos, []modifyLineFunc{
        toGolangReceiverName,
        lowerArgsFirstChar,
        renameInputArgs,
        namedReturnFunc("ans"),
    }, customFuncContent)

    if p.needMod {
        p.defaultCode = "const mod int = 1e9 + 7\n\n" + p.defaultCode
    }

    if err := p.createDir(); err != nil {
        fmt.Println("createDir err:", p.url, err)
        return nil
    }
    if err := p.writeMainFile(); err != nil {
        fmt.Println("writeMainFile err:", p.url, err)
    }
    if err := p.writeTestFile(); err != nil {
        fmt.Println("writeTestFile err:", p.url, err)
    }
    if err := p.writeTestDataFile(); err != nil {
        fmt.Println("writeTestFile err:", p.url, err)
    }
    return nil
}

func handleProblem(session *grequests.Session, p *problem) error {
    if p.openURL {
        if err := open.Run(p.url); err != nil {
            fmt.Println("open err:", p.url, err)
        }
    }

    fmt.Println(p.id, p.url)

    err := parseProblem(session, p)
    return err
}

func parseProblem(session *grequests.Session, p *problem) error {
    if err := p.parseGql(session); err != nil {
        fmt.Println(err)
        return nil
    }

    customFuncContent := "\t\n\treturn"
    if p.needMod {
        customFuncContent = "\t\n\t\n\t\n\tans = (ans%mod + mod) % mod\n\treturn"
    }

    p.defaultCode = modifyDefaultCode(p.defaultCode, p.funcLos, []modifyLineFunc{
        toGolangReceiverName,
        lowerArgsFirstChar,
        renameInputArgs,
        namedReturnFunc("ans"),
    }, customFuncContent)

    if p.needMod {
        p.defaultCode = "const mod int = 1e9 + 7\n\n" + p.defaultCode
    }

    if err := p.createDir(); err != nil {
        fmt.Println("createDir err:", p.url, err)
        return nil
    }
    if err := p.writeMainFile(); err != nil {
        fmt.Println("writeMainFile err:", p.url, err)
    }
    if err := p.writeTestFile(); err != nil {
        fmt.Println("writeTestFile err:", p.url, err)
    }
    if err := p.writeTestDataFile(); err != nil {
        fmt.Println("writeTestFile err:", p.url, err)
    }
    return nil
}

func GenLeetCodeTests(
    username string,
    password string,
    contestTag string,
    openWebPage bool,
    contestDir string,
    customComment string) error {

    session, err := login(username, password)
    if err != nil {
        return err
    }
    fmt.Println("登录成功", host, username)

    var problems []*problem
    for {
        problems, err = fetchProblemURLs(session, contestTag)
        if err == nil {
            break
        }
        fmt.Println(err)
        time.Sleep(500 * time.Millisecond)
    }

    if len(problems) == 0 {
        return nil
    }

    customComment = updateComment(customComment)

    for _, p := range problems {
        p.openURL = openWebPage
        p.customComment = customComment
        p.contestDir = contestDir
    }

    fmt.Println("题目链接获取成功，开始解析")
    return handleContestProblems(session, problems)
}

func GenLeetCodeSingleTest(
    username string,
    password string,
    problemId int,
    openWebPage bool,
    problemDir string,
    customComment string) (err error) {

    session, err := login(username, password)
    if err != nil {
        return err
    }
    fmt.Println("登录成功", host, username)

    var p *problem

    for {
        p, err = fetchSingleProblemURL(session, problemId)
        if err == nil {
            break
        }
        fmt.Println(err)
        time.Sleep(500 * time.Millisecond)
    }

    customComment = updateComment(customComment)

    p.openURL = openWebPage
    p.customComment = customComment
    p.problemDir = problemDir

    fmt.Println("题目链接获取成功，开始解析")
    return handleProblem(session, p)
}
