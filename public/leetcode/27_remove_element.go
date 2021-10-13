package main

import "fmt"

func removeElement(nums []int, val int) int {
	// 默认情况下 Golang 的数组传递是值传递.
	// 想要对传入的数组进行修改时，可以使用指针来对数组进行操作
	// Golang 中的切片与 Java 中的 ArrayList 集合类似，进行的是引用传递
	for i := 0; i < len(nums); i++ {
		if nums[i] == val {
			nums = append(nums[:i], nums[i+1:]...)
			i--
		}
	}
	return len(nums)
}

func main() {
	//make构造一个切片
	nums := []int{3,2,2,3}
	val := 3
	fmt.Println(removeElement(nums, val))
	fmt.Println(nums)
	fmt.Println(len(nums))
}
