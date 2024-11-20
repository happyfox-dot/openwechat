package main

import (
    "fmt"
    "strings"
)

func parseString(input string) (map[string]string, error) {
    result := make(map[string]string)
    // 修正这里的FieldsFunc调用，规范匿名函数写法
    pairs := strings.FieldsFunc(input, func(r rune) bool {
        return r == rune(' ')
    })
    for _, pair := range pairs {
        parts := strings.SplitN(pair, ":", 2)
        if len(parts)!= 2 {
            return nil, fmt.Errorf("invalid pair: %s", pair)
        }
        result[parts[0]] = parts[1]
    }
    return result, nil
}

func main() {
    input := "cmd:2k type:3_pix_gaussian"
    parsed, err := parseString(input)
    if err!= nil {
        fmt.Println("Error parsing string:", err)
        return
    }
    fmt.Println("Parsed result:")
    fmt.Println("Command:", parsed["cmd"])
    fmt.Println("Type:", parsed["type"])
}