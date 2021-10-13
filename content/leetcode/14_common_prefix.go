package main

import "fmt"

func longestCommonPrefix(strs []string) string {
	result := ""

	if len(strs) == 1 {
		return strs[0]
	}

	minStr := ""
	for i, v := range strs {
		if i == 0 {
			minStr = v
		}
		if len(minStr) > len(v) {
			minStr = v
		}
	}

	difFlag := false
	for i := 0; i < len(minStr); i++ {
		for _, v := range strs {
			if v == minStr {
				continue
			}
			difFlag = true
			if minStr[:len(minStr)-i] != v[:len(minStr)-i] {
				result = ""
				break
			}
			result = minStr[:len(minStr)-i]
		}
		if result != "" {
			return result
		}
	}
	if difFlag {
		return result
	} else {
		return minStr
	}
}

func main() {
	strs := []string{"flower","flower","flower","flower"}
	fmt.Println(longestCommonPrefix(strs))
}