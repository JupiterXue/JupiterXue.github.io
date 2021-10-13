package main

import (
	"fmt"
	"strconv"
)

func isPalindrome(x int) bool {
	xStr :=  strconv.Itoa(x)
	halfLen := len(xStr)/2
	wholeLen := len(xStr)
	for i, v := range xStr[:halfLen+1] {
		if int(v) == int(xStr[wholeLen-1-i]) {
			continue
		} else {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(isPalindrome(-121))
}


//fmt.Println(xStr)
//for i := 0; i < len(xStr) +1 ; i++ {
//value := xStr[i]
//reverse_value := xStr[len(xStr)-1-i]
//fmt.Println(i)
//fmt.Println(value)
//fmt.Println(reverse_value)
//fmt.Println()
//if int(value) == int(reverse_value) {
//continue
//} else {
//return false
//}
//}
//return true