package main

import "fmt"

func smallerNumbersThanCurrent(nums []int) []int {
	result := []int{}
	for i, v1 := range nums {
		num := 0
		for j, v2 := range nums {
			if i == j {
				continue
			}
			if v1 > v2 {
				num += 1
			}
		}
		result = append(result, num)
	}
	return result
}

func main() {
	nums := []int{8,1,2,2,3}
	fmt.Println(smallerNumbersThanCurrent(nums))
}
