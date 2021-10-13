package main

func twoSum(nums []int, target int) []int {
	result := []int{}
	for index, v := range  nums {
		//if v <= target {
		for index2, v2 := range  nums {
			if index == index2 {
				continue
			} else {
				if v + v2 == target {
					result = append(result, index)
					result = append(result, index2)
					return result
				}
			}
		}
		//}
	}
	return result
}

func main() {
	//nums := []int{2,7,11,15}
	//target := 9

	nums := []int{0,4,3,0}
	target := 0

	twoSum(nums, target)
}
