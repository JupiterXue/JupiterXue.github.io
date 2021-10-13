package main

import (
	"fmt"
	"math"
)

func is_out(x int) bool {
	// 如果反转后整数超过 32 位的有符号整数的范围 [−2^31,  2^31 − 1] ，就返回 0。
	top := math.Pow(float64(2), float64(31)) -1
	bottom := math.Pow(float64(-2), float64(31))
	if bottom <= float64(x) && float64(x) <= top {
		return false
	} else {
		return true
	}
}

func reverse(x int) int {
	// 把每一位取出来，装入数组。遍历数组，判断大小，如果大于 0 则乘 10，否则乘 -10。
	// Golang 除法结果为向下取得整数
	result_list := []int{}
	for i := int(math.Abs(float64(x))); i > 0 ; i = i/10 {
		result_list = append(result_list, i%10)
	}
	result_int := 0
	for i, v := range result_list{
		flat := math.Pow(float64(10), float64(len(result_list) - 1 - i))
		result_int += int(flat) * v
	}

	if is_out(result_int) {
		return 0
	} else {
		if x > 0 {
			return result_int
		} else {
			return result_int * -1
		}
	}
}

func main() {
	fmt.Println(reverse(1534236469))
}