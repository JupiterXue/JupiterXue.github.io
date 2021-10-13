package main

import "fmt"

func romanToInt(s string) int {
	dic := map[string]int{
		"I": 1,
		"V": 5,
		"X": 10,
		"L": 50,
		"C": 100,
		"D": 500,
		"M": 1000,
	} //表示创建一个key为string，value的值为int的数据类型。
	specDic := map[string]int{
		"IV": 4,
		"IX": 9,
		"XL": 40,
		"XC": 90,
		"CD": 400,
		"CM": 900,
	}

	specIndexDic := map[int]int{}

	for i := 0; i < len(s)-1 ; i++ {
		if _, ok := specDic[s[i:i+2]]; ok {
			//存在
			specIndexDic[i] = specDic[s[i:i+2]]
		}
	}
	//fmt.Println(specIndexDic)

	var valueList []int
	maxValue := 0
	maxValueIndex := 0
	for i, v := range s {
		currentValue := dic[string(v)]
		if _, ok := specIndexDic[i-1]; ok {
			continue
		} else if _, ok := specIndexDic[i]; ok {
			currentValue = specIndexDic[i]
			valueList = append(valueList, currentValue)
		} else {
			valueList = append(valueList, currentValue)
		}
		if maxValue < currentValue { // 记录最大值索引位置
			maxValueIndex = len(valueList) - 1
			maxValue = currentValue
		}
	}
	fmt.Println(valueList)
	removeValue := 0
	plusValue := 0
	for i, v := range valueList {
		if i < maxValueIndex {
			removeValue += v
		} else if i > maxValueIndex {
			plusValue += v
		}
	}
	result := valueList[maxValueIndex]
	result -= removeValue
	result += plusValue
	return result
}

func main() {
	r := romanToInt("MMMXLV")
	// attempt = 3045
	// result = 3065
	fmt.Println(r)
}