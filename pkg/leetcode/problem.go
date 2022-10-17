package leetcode

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/levigross/grequests"
    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
    "io"
    "os"
    "strconv"
    "strings"
    "time"
)

type problem struct {
    id      string
    url     string
    openURL bool

    defaultCode   string
    funcName      string
    isFuncProblem bool
    needMod       bool
    funcLos       []int
    customComment string

    sampleIns  [][]string
    sampleOuts [][]string

    contestDir string
    problemDir string
}

// 解析一个样例输入或输出
func (p *problem) parseSampleText(text string, parseArgs bool) []string {
    text = strings.ReplaceAll(text, " ", " ") // 替换 NBSP 为正常空格

    text = strings.TrimSpace(text)
    if text == "" {
        return nil
    }

    lines := strings.Split(text, "\n")
    for i, s := range lines {
        lines[i] = strings.TrimSpace(s)
    }

    // 由于新版的样例不是这种格式了，这种特殊情况就不处理了
    // 见 https://leetcode-cn.com/contest/weekly-contest-121/problems/time-based-key-value-store/
    if !p.isFuncProblem {
        return lines
    }

    text = strings.Join(lines, "")

    // 包含中文的话，说明原始数据有误，截断首个中文字符之后的字符
    if idx := findNonASCII(text); idx != -1 {
        fmt.Println("[warn] 样例数据含有非 ASCII 字符，截断，原文为", text)
        text = text[:idx]
    }

    // 不含等号，说明只有一个参数
    if !parseArgs || !strings.Contains(text, "=") {
        return []string{text}
    }

    // TODO: 处理参数本身含有 = 的情况
    splits := strings.Split(text, "=")
    sample := make([]string, 0, len(splits)-1)
    for _, s := range splits[1 : len(splits)-1] {
        end := strings.LastIndexByte(s, ',')
        sample = append(sample, strings.TrimSpace(s[:end]))
    }
    sample = append(sample, strings.TrimSpace(splits[len(splits)-1]))
    if !p.isFuncProblem {
        sample = []string{strings.Join(sample, "\n") + "\n"}
    }
    return sample
}

func (p *problem) parsePossibleSampleTexts(texts []string, parseArgs bool) []string {
    for _, text := range texts {
        if sample := p.parseSampleText(text, parseArgs); len(sample) > 0 {
            return sample
        }
    }
    return nil
}

func (p *problem) parseGql(session *grequests.Session) error {
    var slug string
    strs := strings.Split(p.url, "/")
    for i := len(strs) - 1; i >= 0; i-- {
        if strs[i] != "" {
            slug = strs[i]
            break
        }
    }
    resp, err := session.Post(graphqlURL, &grequests.RequestOptions{
        JSON: map[string]any{
            "operationName": "questionData",
            "variables": map[string]any{
                "titleSlug": slug,
            },
            "query": "query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionFrontendId\n    categoryTitle\n    boundTopicId\n    title\n    titleSlug\n    content\n    translatedTitle\n    translatedContent\n    isPaidOnly\n    difficulty\n    likes\n    dislikes\n    isLiked\n    similarQuestions\n    contributors {\n      username\n      profileUrl\n      avatarUrl\n      __typename\n    }\n    langToValidPlayground\n    topicTags {\n      name\n      slug\n      translatedName\n      __typename\n    }\n    companyTagStats\n    codeSnippets {\n      lang\n      langSlug\n      code\n      __typename\n    }\n    stats\n    hints\n    solution {\n      id\n      canSeeDetail\n      __typename\n    }\n    status\n    sampleTestCase\n    metaData\n    judgerAvailable\n    judgeType\n    mysqlSchemas\n    enableRunCode\n    envInfo\n    book {\n      id\n      bookName\n      pressName\n      source\n      shortDescription\n      fullDescription\n      bookImgUrl\n      pressImgUrl\n      productUrl\n      __typename\n    }\n    isSubscribed\n    isDailyQuestion\n    dailyRecordStatus\n    editorType\n    ugcQuestionId\n    style\n    exampleTestcases\n    jsonExampleTestcases\n    __typename\n  }\n}\n",
        },
    })
    if err != nil {
        fmt.Println("获取题目详情失败，重试", err)
        time.Sleep(time.Second)
        return p.parseGql(session)
    }
    if !resp.Ok {
        return fmt.Errorf("POST %s return code %d", graphqlURL, resp.StatusCode)
    }

    d := struct {
        Data struct {
            Question struct {
                QuestionId        string `json:"questionFrontendId"`
                TitleSlug         string `json:"titleSlug"`
                TranslatedTitle   string `json:"translatedTitle"`
                TranslatedContent string `json:"translatedContent"`
                CodeSnippets      []struct {
                    LangSlug string `json:"langSlug"`
                    Code     string `json:"code"`
                } `json:"codeSnippets"`
            } `json:"question"`
        } `json:"data"`
    }{}
    if err = resp.JSON(&d); err != nil {
        return err
    }

    question := d.Data.Question
    p.needMod = strings.Contains(question.TranslatedContent, "取余")
    p.needMod = p.needMod || strings.Contains(question.TranslatedContent, "取模")
    p.needMod = p.needMod || strings.Contains(question.TranslatedContent, "答案可能很大")

    for _, snippet := range question.CodeSnippets {
        if snippet.LangSlug == "golang" {
            p.defaultCode = strings.TrimSpace(snippet.Code)
            // 下面解析样例需要知道 p.isFuncProblem
            p.funcName, p.isFuncProblem, p.funcLos = parseCode(p.defaultCode)
            break
        }
    }
    if p.defaultCode == "" {
        fmt.Println("解析失败，未找到 Go 代码模板！")
    }

    lines := strings.Split(d.Data.Question.TranslatedContent, "\n")
    for i := 0; i < len(lines); i++ {
        labels := []string{"<strong>输入：</strong>", "<strong>输入: </strong>"}
        for _, label := range labels {
            if strings.Contains(lines[i], label) {
                pos := strings.Index(lines[i], label)
                p.sampleIns = append(p.sampleIns, p.parseSampleText(lines[i][pos+len(label):], true))
                break
            }
        }

        labels = []string{"<strong>输出：</strong>", "<strong>输出: </strong>"}
        for _, label := range labels {
            if strings.Contains(lines[i], label) {
                pos := strings.Index(lines[i], label)
                p.sampleOuts = append(p.sampleOuts, p.parseSampleText(lines[i][pos+len(label):], true))
                break
            }
        }
    }
    if len(p.sampleIns) == 0 {
        return fmt.Errorf("解析失败，未找到样例输入输出！")
    }
    if len(p.sampleIns) != len(p.sampleOuts) {
        return fmt.Errorf("len(sampleIns) != len(sampleOuts) : %d != %d", len(p.sampleIns), len(p.sampleOuts))
    }
    return nil
}

// 获取题目样例和代码
func (p *problem) parseHTML(session *grequests.Session) (err error) {
    defer func() {
        // visit htmlNode may cause panic
        if er := recover(); er != nil {
            err = fmt.Errorf("need fix: %v", er)
        }
    }()

    resp, err := session.Get(p.url, nil)
    if err != nil {
        return
    }
    if !resp.Ok {
        return fmt.Errorf("GET %s return code %d", p.url, resp.StatusCode)
    }

    htmlText, _ := io.ReadAll(resp)
    p.needMod = bytes.Contains(htmlText, []byte("取余"))
    p.needMod = p.needMod || bytes.Contains(htmlText, []byte("取模"))
    p.needMod = p.needMod || bytes.Contains(htmlText, []byte("答案可能很大"))

    rootNode, err := html.Parse(bytes.NewReader(htmlText))
    if err != nil {
        return err
    }

    htmlNode := rootNode.FirstChild.NextSibling
    var bodyNode *html.Node
    for o := htmlNode.FirstChild; o != nil; o = o.NextSibling {
        if o.DataAtom == atom.Body {
            bodyNode = o
            break
        }
    }

    // parse defaultCode
    for o := bodyNode.FirstChild; o != nil; o = o.NextSibling {
        if o.DataAtom == atom.Script && o.FirstChild != nil {
            jsText := o.FirstChild.Data
            if start := strings.Index(jsText, "codeDefinition:"); start != -1 {
                end := strings.Index(jsText, "enableTestMode")
                jsonText := jsText[start+len("codeDefinition:") : end]
                jsonText = strings.TrimSpace(jsonText)
                jsonText = jsonText[:len(jsonText)-3] + "]" // remove , at end
                jsonText = strings.Replace(jsonText, `'`, `"`, -1)

                var data []struct {
                    Value       string `json:"value"`
                    DefaultCode string `json:"defaultCode"`
                }
                if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
                    return err
                }

                for _, template := range data {
                    if template.Value == "golang" {
                        p.defaultCode = strings.TrimSpace(template.DefaultCode)
                        // 下面解析样例需要知道 p.isFuncProblem
                        p.funcName, p.isFuncProblem, p.funcLos = parseCode(p.defaultCode)
                        break
                    }
                }
                break
            }
        }
    }
    if p.defaultCode == "" {
        fmt.Println("解析失败，未找到 Go 代码模板！")
    }

    var parseNode func(*html.Node)
    parseNode = func(o *html.Node) {
        // 提取并解析每个 <pre> 块内的文本（以中文为基准解析）
        // 需要判断 <pre> 的下一个子元素是否为 tag
        //     https://leetcode-cn.com/contest/weekly-contest-190/problems/max-dot-product-of-two-subsequences/
        //     https://leetcode-cn.com/contest/weekly-contest-212/problems/arithmetic-subarrays/
        // 有 tag 也不一定为 <strong>
        //     <img> https://leetcode-cn.com/contest/weekly-contest-103/problems/snakes-and-ladders/
        //     <b> https://leetcode-cn.com/contest/weekly-contest-210/problems/split-two-strings-to-make-palindrome/
        //     <code> https://leetcode-cn.com/contest/weekly-contest-163/problems/shift-2d-grid/
        // 提取出文本后，去掉「解释」和「提示」后面的文字，然后分「输入」和「输出」来解析后面的数据
        if o.DataAtom == atom.Pre && o.FirstChild.DataAtom != 0 && o.FirstChild.DataAtom != atom.Img && o.FirstChild.DataAtom != atom.Image { // 一般是 atom.Strong，特殊情况是 atom.B
            // 找到第一个文本，这样写是因为可能有额外的嵌套 tag https://leetcode-cn.com/contest/weekly-contest-163/problems/shift-2d-grid/
            var data string
            for o := o.FirstChild.FirstChild; o != nil; o = o.FirstChild {
                if o.DataAtom == 0 {
                    data = o.Data
                    break
                }
            }
            if strings.HasPrefix(strings.TrimSpace(data), "输") { // 输入（极少情况下会被错误地写成输出）
                rawData := &strings.Builder{}
                var parsePreNode func(*html.Node)
                parsePreNode = func(o *html.Node) {
                    if o.DataAtom == 0 {
                        rawData.WriteString(o.Data)
                    }
                    for c := o.FirstChild; c != nil; c = c.NextSibling {
                        parsePreNode(c)
                    }
                }
                parsePreNode(o)

                data := rawData.String()
                if i := strings.Index(data, "解"); i >= 0 { // 解释
                    data = data[:i]
                }
                if i := strings.Index(data, "提"); i >= 0 { // 提示
                    data = data[:i]
                }

                // 去掉前两个字和冒号
                data = strings.TrimSpace(data)
                if i := strings.IndexRune(data, ':'); i >= 0 {
                    data = data[i+1:]
                } else if i := strings.IndexRune(data, '：'); i >= 0 {
                    data = data[i+3:]
                }

                i := strings.Index(data, "输") // 输出
                p.sampleIns = append(p.sampleIns, p.parseSampleText(data[:i], true))

                // 去掉前两个字和冒号
                if i := strings.IndexRune(data, ':'); i >= 0 {
                    data = data[i+1:]
                } else if i := strings.IndexRune(data, '：'); i >= 0 {
                    data = data[i+3:]
                }
                p.sampleOuts = append(p.sampleOuts, p.parseSampleText(data, true))
                return
            }
        }
        for c := o.FirstChild; c != nil; c = c.NextSibling {
            parseNode(c)
        }
    }
    parseNode(bodyNode)

    if len(p.sampleIns) == 0 {
        // 没找到 <pre>，国服特殊比赛（春秋赛等）
        parseNode = func(o *html.Node) {
            if o.DataAtom == atom.Div && o.FirstChild != nil && strings.Contains(o.FirstChild.Data, "示例") {
                raw := o.FirstChild.Data
                sp := strings.Split(raw, "`")
                for i, s := range sp {
                    if strings.Contains(s, ">输入") || strings.Contains(s, "> 输入") {
                        text := sp[i+1]
                        if !p.isFuncProblem {
                            // https://leetcode-cn.com/contest/season/2020-fall/problems/IQvJ9i/
                            text += "\n" + sp[i+3] // 跳过 sp[i+2]
                        }
                        p.sampleIns = append(p.sampleIns, p.parseSampleText(text, true))
                    } else if strings.Contains(s, ">输出") || strings.Contains(s, "> 输出") {
                        p.sampleOuts = append(p.sampleOuts, p.parseSampleText(sp[i+1], true))
                    }
                }
            }
            for c := o.FirstChild; c != nil; c = c.NextSibling {
                parseNode(c)
            }
        }
        parseNode(bodyNode)
    }

    if len(p.sampleIns) != len(p.sampleOuts) {
        return fmt.Errorf("len(sampleIns) != len(sampleOuts) : %d != %d", len(p.sampleIns), len(p.sampleOuts))
    }
    if len(p.sampleIns) == 0 {
        return fmt.Errorf("解析失败，未找到样例输入输出！")
    }
    return nil
}

func (p *problem) createDir() error {
    if p.problemDir != "" {
        return os.MkdirAll(p.problemDir, os.ModePerm)
    }
    return os.MkdirAll(p.contestDir+p.id, os.ModePerm)
}

func (p *problem) writeMainFile() error {
    imports := ""
    if strings.Contains(p.defaultCode, "Definition for") {
        imports = `
import . "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
`
    }
    p.defaultCode = strings.TrimSpace(p.defaultCode)
    fileContent := fmt.Sprintf(`package main
%s
%s%s
`, imports, p.customComment, p.defaultCode)

    filePath := p.contestDir + fmt.Sprintf("%[1]s/%[1]s.go", p.id)
    if p.problemDir != "" {
        filePath = p.problemDir + "solution.go"
    }
    return os.WriteFile(filePath, []byte(fileContent), 0644)
}

func (p *problem) writeTestFile() error {
    logInfo := ""
    testUtilFunc := "testutils.RunLeetCodeFuncWithFile"
    if !p.isFuncProblem {
        logInfo += "\n\t" + `t.Log("记得初始化所有全局变量")`
        testUtilFunc = "testutils.RunLeetCodeClassWithFile"
    }
    var testName, testDataFileName string
    if p.problemDir != "" {
        s, _ := strconv.Atoi(p.id)
        testName = fmt.Sprintf("%04d", s)
        testDataFileName = "sample.txt"
    } else {
        testName = strings.ToUpper(p.id)
        testDataFileName = fmt.Sprintf("%s.txt", p.id)
    }
    testStr := fmt.Sprintf(`package main

import (
    "github.com/changaolee/go-algorithms/pkg/leetcode/testutils"
    "testing"
)

func Test%s(t *testing.T) {%s
    targetCaseNum := 0 // -1
    if err := %s(t, %s, "%s", targetCaseNum); err != nil {
        t.Fatal(err)
    }
}

// %s
`, testName, logInfo, testUtilFunc, p.funcName, testDataFileName, p.url)

    filePath := p.contestDir + fmt.Sprintf("%[1]s/%[1]s_test.go", p.id)
    if p.problemDir != "" {
        filePath = p.problemDir + "solution_test.go"
    }
    return os.WriteFile(filePath, []byte(testStr), 0644)
}

func (p *problem) writeTestDataFile() error {
    var lines []string
    for i, inArgs := range p.sampleIns {
        lines = append(lines, inArgs...)
        if i < len(p.sampleOuts) {
            lines = append(lines, p.sampleOuts[i]...)
        }
        lines = append(lines, "") // empty line for clarity
    }
    lines = append(lines, "", "", "")
    testDataStr := strings.Join(lines, "\n")

    filePath := p.contestDir + fmt.Sprintf("%[1]s/%[1]s.txt", p.id)
    if p.problemDir != "" {
        filePath = p.problemDir + "sample.txt"
    }
    return os.WriteFile(filePath, []byte(testDataStr), 0644)
}
