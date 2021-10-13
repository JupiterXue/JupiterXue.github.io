package main

import "fmt"

func strStr(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}
	startIndex := 0
	lenNeedle := len(needle)
	startChar := string(needle[0])
	for i, v := range haystack{
		if string(v) == startChar {
			startIndex = i
			if len(haystack[startIndex:]) >= lenNeedle {
				if haystack[startIndex:startIndex+lenNeedle] == needle {
					return startIndex
				}
			}
		}
	}
	return -1
}

func main() {
	haystack := "hello"
	needle := "ll"
	fmt.Println(strStr(haystack, needle))
}
