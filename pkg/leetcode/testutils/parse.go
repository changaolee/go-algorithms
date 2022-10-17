package testutils

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

func parseRawArg(tp reflect.Type, rawData string) (v reflect.Value, err error) {
    rawData = strings.TrimSpace(rawData)
    invalidErr := fmt.Errorf("invalid test data: %s", rawData)
    switch tp.Kind() {
    case reflect.String:
        if len(rawData) <= 1 || rawData[0] != '"' && rawData[0] != '\'' || rawData[len(rawData)-1] != rawData[0] {
            return reflect.Value{}, invalidErr
        }
        // remove " or ' at leftmost and rightmost
        v = reflect.ValueOf(rawData[1 : len(rawData)-1])
    case reflect.Uint8: // byte
        // rawData like "a" or 'a'
        if len(rawData) != 3 || rawData[0] != '"' && rawData[0] != '\'' || rawData[2] != rawData[0] {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(rawData[1])
    case reflect.Int:
        i, er := strconv.Atoi(rawData)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(i)
    case reflect.Uint:
        i, er := strconv.Atoi(rawData)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(uint(i))
    case reflect.Int64:
        i, er := strconv.ParseInt(rawData, 10, 64)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(i)
    case reflect.Uint64:
        i, er := strconv.ParseUint(rawData, 10, 64)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(i)
    case reflect.Float64:
        f, er := strconv.ParseFloat(rawData, 64)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(f)
    case reflect.Bool:
        b, er := strconv.ParseBool(rawData)
        if er != nil {
            return reflect.Value{}, invalidErr
        }
        v = reflect.ValueOf(b)
    case reflect.Slice:
        splits, er := parseRawArray(rawData)
        if er != nil {
            return reflect.Value{}, er
        }
        v = reflect.New(tp).Elem()
        for _, s := range splits {
            _v, er := parseRawArg(tp.Elem(), s)
            if er != nil {
                return reflect.Value{}, er
            }
            v = reflect.Append(v, _v)
        }
    case reflect.Ptr: // *TreeNode, *ListNode, *Point, *Interval
        switch tpName := tp.Elem().Name(); tpName {
        case "TreeNode":
            root, er := buildTreeNode(rawData)
            if er != nil {
                return reflect.Value{}, er
            }
            v = reflect.ValueOf(root)
        case "ListNode":
            head, er := buildListNode(rawData)
            if er != nil {
                return reflect.Value{}, er
            }
            v = reflect.ValueOf(head)
        case "Point":
            p, er := buildPoint(rawData)
            if er != nil {
                return reflect.Value{}, er
            }
            v = reflect.ValueOf(p)
        case "Interval":
            p, er := buildInterval(rawData)
            if er != nil {
                return reflect.Value{}, er
            }
            v = reflect.ValueOf(p)
        default:
            return reflect.Value{}, fmt.Errorf("unknown type %s", tpName)
        }
    default:
        return reflect.Value{}, fmt.Errorf("unknown type %s", tp.Name())
    }
    return
}

func parseRawArray(rawArray string) (splits []string, err error) {
    invalidErr := fmt.Errorf("invalid test data: %s", rawArray)

    // check [] at leftmost and rightmost
    if len(rawArray) <= 1 || rawArray[0] != '[' || rawArray[len(rawArray)-1] != ']' {
        return nil, invalidErr
    }

    // ignore [] at leftmost and rightmost
    rawArray = rawArray[1 : len(rawArray)-1]
    if rawArray == "" {
        return
    }

    isPoint := rawArray[0] == '('

    const sep = ','
    var depth, quotCnt, bracketCnt int
    for start := 0; start < len(rawArray); {
        end := start
    outer:
        for ; end < len(rawArray); end++ {
            switch rawArray[end] {
            case '[':
                depth++
            case ']':
                depth--
            case '"':
                quotCnt++
            case '(':
                bracketCnt++
            case ')':
                bracketCnt--
            case sep:
                if depth == 0 {
                    if !isPoint {
                        if quotCnt%2 == 0 {
                            break outer
                        }
                    } else {
                        if bracketCnt%2 == 0 {
                            break outer
                        }
                    }
                }
            }
        }
        splits = append(splits, strings.TrimSpace(rawArray[start:end]))
        start = end + 1 // skip sep
    }
    if depth != 0 || quotCnt%2 != 0 {
        return nil, invalidErr
    }
    return
}

func toRawString(v reflect.Value) (s string, err error) {
    switch v.Kind() {
    case reflect.Slice:
        sb := &strings.Builder{}
        sb.WriteByte('[')
        for i := 0; i < v.Len(); i++ {
            if i > 0 {
                sb.WriteByte(',')
            }
            _s, er := toRawString(v.Index(i))
            if er != nil {
                return "", er
            }
            sb.WriteString(_s)
        }
        sb.WriteByte(']')
        s = sb.String()
    case reflect.Ptr: // *TreeNode, *ListNode, *Point, *Interval
        switch tpName := v.Type().Elem().Name(); tpName {
        case "TreeNode":
            s = v.Interface().(*TreeNode).toRawString()
        case "ListNode":
            s = v.Interface().(*ListNode).toRawString()
        case "Point":
            s = v.Interface().(*Point).toRawString()
        case "Interval":
            s = v.Interface().(*Interval).toRawString()
        default:
            return "", fmt.Errorf("unknown type %s", tpName)
        }
    case reflect.String:
        s = fmt.Sprintf(`"%s"`, v.Interface())
    case reflect.Uint8: // byte
        s = fmt.Sprintf(`"%c"`, v.Interface())
    case reflect.Float64:
        s = fmt.Sprintf(`%.5f`, v.Interface())
    default: // int uint int64 uint64 bool
        s = fmt.Sprintf(`%v`, v.Interface())
    }
    return
}
