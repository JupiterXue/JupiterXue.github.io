package main

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	result := new(ListNode)
	point := result
	sum := 0
	for {
		if l1 != nil && l2 != nil{
		 	sum += l1.Val + l2.Val
			l1 = l1.Next
			l2 = l2.Next
		} else {
			if l1 != nil {
				sum += l1.Val
				l1 = l1.Next
			}

			if l2 != nil {
				sum += l2.Val
				l2 = l2.Next
			}
		}
		point.Val = sum % 10
		sum /= 10
		if l1 == nil && l2 == nil && sum == 0 {
			break
		} else {
			point.Next = new(ListNode)
			point = point.Next
		}
	}
	return result
}

func PrintlistNode(l0 *ListNode, name string) {
	var result []int
	node := l0
	for ;; {
		if node == nil {
			break
		}
		result = append(result, node.Val)
		node = node.Next
	}
	fmt.Print(name + "=")
	fmt.Println(result)
}

func createListNode(str string) *ListNode {
	result := new(ListNode)
	point := result
	dataList := strings.Split(str, ",")
	for _, v := range dataList {
		value, _ := strconv.Atoi(v)
		tmpNode := &ListNode{
			value,
			nil,
		}
		point.Next = tmpNode
		point = point.Next
	}
	return result.Next
}

func main() {
	//l1 := createListNode("2,4,3")
	//l2 := createListNode("5,6,4")

	l1 := createListNode("9,9,9,9,9,9,9")
	l2 := createListNode("9,9,9,9")

	PrintlistNode(l1, "l1")
	PrintlistNode(l2, "l2")

	PrintlistNode(addTwoNumbers(l1, l2), "result")
	PrintlistNode(l1, "l1")
	PrintlistNode(l2, "l2")
}