package main

//type ListNode struct {
//	Val  int
//	Next *ListNode
//}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	newList := &ListNode{}
	current := newList
	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			current.Next = l1
			l1 = l1.Next
		} else {
			current.Next = l2
			l2 = l2.Next
		}
		current = current.Next
	}
	if l1 == nil {
		current.Next = l2
	} else {
		current.Next = l1
	}
	return newList.Next
}


//func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
//	result := new(ListNode)
//	result = l1
//	node := new(ListNode)
//	node = l1
//	for ;; {
//		if node == nil {
//			node = l2
//			result.Next = node
//			break
//		}
//		node = node.Next
//	}
//	return result
//}

//func PrintlistNode(l0 *ListNode, name string) {
//	var result []int
//	node := l0
//	for ;; {
//		if node == nil {
//			break
//		}
//		result = append(result, node.Val)
//		node = node.Next
//	}
//	fmt.Print(name + "=")
//	fmt.Println(result)
//}

func main() {
	l1 := new(ListNode)
	l1.Val = 1
	next := new(ListNode)
	next.Val = 2
	l1.Next = next

	next2 := new(ListNode)
	next2.Val = 4
	next.Next = next2



	l2 := new(ListNode)
	l2.Val = 1
	next21 := new(ListNode)
	next21.Val = 3
	l2.Next = next21

	next22 := new(ListNode)
	next22.Val = 4
	next21.Next = next22

	//PrintlistNode(l1, "l1")
	//PrintlistNode(l2, "l2")

	PrintlistNode(mergeTwoLists(l1, l2), "result")
	PrintlistNode(l1, "l1")
	PrintlistNode(l2, "l2")
}
