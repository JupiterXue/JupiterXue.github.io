package main

import "fmt"

func isValid(s string) bool {
	bracketsDic := map[string]string{
		"(" : ")",
		"[" : "]",
		"{" : "}",
	}

	var stack []string
	stack = append(stack, string(s[0]))
	for i := 1; i < len(s); i++ {
		if len(stack) == 0 {
			stack = append(stack, string(s[i]))
			continue
		}
		if _, ok := bracketsDic[string(s[i])]; ok {
			stack = append(stack, string(s[i]))
			continue
		}
			preBracket := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if bracketsDic[preBracket] == string(s[i]) {
			continue
		}
		return false
	}
	if len(stack) != 0 {
		return false
	} else {
		return true
	}
}

func main() {
	s := "{[]}"
	fmt.Println(isValid(s))
}