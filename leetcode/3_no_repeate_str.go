package main

import "fmt"

func check(partS string) bool {
	for i = 0; i < len(partS) - 1; i++ {
		if partS[i] == partS[i+1] {
			return false
		}
	}
	return true
}

func lengthOfLongestSubstring(s string) int {
	lenOrder := 0
	for i = 0; i < len(s); i++ {
		if check(s[i:]) {
			lenOrder = len(s[i:])
			break
		}
	}

	lenROrder := 0
	for i = 0; i < len(s); i++ {
		if check(s[:len(s)-i]) {
			lenROrder  = len(s[:len(s)-i])
		}
	}

	if lenOrder >= lenROrder {
		return lenOrder
	} else {
		return lenOrder
	}
}

func main() {
	s := "abcabcbb"
	fmt.Println(lengthOfLongestSubstring(s))
}
