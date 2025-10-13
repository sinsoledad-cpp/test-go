package main

import "fmt"

func myAppend(s []int) []int {
	s = append(s, 100)
	return s
}

func myAppendPtr(s *[]int) {
	*s = append(*s, 100)
	return
}

func main() {
	s := []int{1, 1, 1}
	newS := myAppend(s)

	fmt.Println(s)
	fmt.Println(newS)

	s = newS

	myAppendPtr(&s)
	fmt.Println(s)
}
