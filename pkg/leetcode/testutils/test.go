package testutils

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "os"
    "reflect"
    "strings"
    "testing"
)

func RunLeetCodeFuncWithFile(t *testing.T, f interface{}, filePath string, targetCaseNum int) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    lines := trimSpaceAndEmptyLine(string(data))
    n := len(lines)
    if n == 0 {
        return fmt.Errorf("输入数据为空，请检查文件路径和文件内容是否正确")
    }

    fType := reflect.TypeOf(f)
    if fType.Kind() != reflect.Func {
        return fmt.Errorf("f 必须是函数")
    }

    // 每 fNumIn+fNumOut 行一组数据
    fNumIn := fType.NumIn()
    fNumOut := fType.NumOut()
    tcSize := fNumIn + fNumOut
    if n%tcSize != 0 {
        return fmt.Errorf("有效行数 %d，应该是 %d 的倍数", n, tcSize)
    }

    examples := make([][]string, 0, n/tcSize)
    for i := 0; i < n; i += tcSize {
        examples = append(examples, lines[i:i+tcSize])
    }
    return RunLeetCodeFuncWithExamples(t, f, examples, targetCaseNum)
}

// RunLeetCodeFuncWithExamples
// rawExamples[i] = 输入+输出，若反射出来的函数或 rawExamples 数据不合法，则会返回一个非空的 error，否则返回 nil
func RunLeetCodeFuncWithExamples(t *testing.T, f interface{}, rawExamples [][]string, targetCaseNum int) (err error) {
    fType := reflect.TypeOf(f)
    if fType.Kind() != reflect.Func {
        return fmt.Errorf("f must be a function")
    }

    fNumIn := fType.NumIn()
    fNumOut := fType.NumOut()

    // 例如，-1 表示最后一个测试用例
    if targetCaseNum < 0 {
        targetCaseNum += len(rawExamples) + 1
    }

    allCasesOk := true
    fValue := reflect.ValueOf(f)
    for curCaseNum, example := range rawExamples {
        if targetCaseNum > 0 && curCaseNum+1 != targetCaseNum {
            continue
        }

        if len(example) != fNumIn+fNumOut {
            return fmt.Errorf("len(example) = %d, but we need %d+%d", len(example), fNumIn, fNumOut)
        }

        rawIn := example[:fNumIn]
        ins := make([]reflect.Value, len(rawIn))
        for i, rawArg := range rawIn {
            rawArg = trimSpace(rawArg)
            ins[i], err = parseRawArg(fType.In(i), rawArg)
            if err != nil {
                return
            }
        }
        // just check rawExpectedOuts is valid or not
        rawExpectedOuts := example[fNumIn:]
        for i := range rawExpectedOuts {
            rawExpectedOuts[i] = trimSpace(rawExpectedOuts[i])
            if _, err = parseRawArg(fType.Out(i), rawExpectedOuts[i]); err != nil {
                return
            }
        }

        const maxInputSize = 150
        inputInfo := strings.Join(rawIn, "\n")
        if len(inputInfo) > maxInputSize { // 截断过长的输入
            inputInfo = inputInfo[:maxInputSize] + "..."
        }

        var outs []reflect.Value
        _f := func() { outs = fValue.Call(ins) }
        if targetCaseNum == 0 && isTLE(_f) {
            allCasesOk = false
            t.Errorf("Time Limit Exceeded %d\nInput:\n%s", curCaseNum+1, inputInfo)
            continue
        } else if targetCaseNum != 0 {
            _f()
        }

        t.Run(fmt.Sprintf("Case %d", curCaseNum+1), func(t *testing.T) {
            for i, out := range outs {
                rawActualOut, er := toRawString(out)
                if er != nil {
                    t.Fatal(er)
                }
                if AssertOutput && !assert.Equal(t, rawExpectedOuts[i], rawActualOut, "Wrong Answer %d\nInput:\n%s", curCaseNum+1, inputInfo) {
                    allCasesOk = false
                }
            }
        })
    }

    // 若有测试用例未通过，则前面必然会打印一些信息，这里直接返回
    if !allCasesOk {
        return nil
    }

    // 若测试的是单个用例，则接着测试所有用例
    if targetCaseNum > 0 {
        t.Logf("case %d is passed", targetCaseNum)
        return RunLeetCodeFuncWithExamples(t, f, rawExamples, 0)
    }

    t.Log("OK")
    return nil
}
