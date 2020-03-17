//https://www.openmymind.net/assets/go/go.pdf
package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Saiyan struct {
	Name  string
	Power int
}

func main() {
	// var goku *Saiyan
	// goku = &Saiyan{"Goku", 9000}
	goku := &Saiyan{"Goku", 9000}
	fmt.Printf("g, %v, %p\n", goku.Power, goku)
	Super(goku)
	fmt.Printf("g, %v, %p\n", goku.Power, goku)

	pointers()
	slices()
}

func byval(q *int) {
	fmt.Printf("3. byval -- q %T: &q=%p q=&i=%p  *q=i=%v\n", q, &q, q, *q)
	*q = 4143
	fmt.Printf("4. byval -- q %T: &q=%p q=&i=%p  *q=i=%v\n", q, &q, q, *q)
	q = nil
}

func pointers() {
	i := int(42)
	fmt.Printf("1. main  -- i  %T: &i=%p i=%v\n", i, &i, i)
	p := &i
	fmt.Printf("2. main  -- p %T: &p=%p p=&i=%p  *p=i=%v\n", p, &p, p, *p)
	byval(p)
	fmt.Printf("5. main  -- p %T: &p=%p p=&i=%p  *p=i=%v\n", p, &p, p, *p)
	fmt.Printf("6. main  -- i  %T: &i=%p i=%v\n", i, &i, i)
}

func Super(s *Saiyan) {
	s.Power = 3000
	fmt.Printf("s, %v, %p\n", s.Power, s)
	s = &Saiyan{"Gohan", 1000}
	fmt.Printf("s, %v, %p\n", s.Power, s)
}

func removeAtIndex(source []int, index int) []int {
	lastIndex := len(source) - 1
	//swap the last value and the value we want to remove
	source[index], source[lastIndex] = source[lastIndex], source[index]
	return source[:lastIndex]
}

func slices() {
	scores := make([]int, 5)
	scores = append(scores, 9332)
	fmt.Println(scores)

	scores1 := []int{1, 2, 3, 4, 5}
	scores1 = removeAtIndex(scores1, 2)
	fmt.Println(scores1)

	scores2 := make([]int, 100)
	for i := 0; i < 100; i++ {
		scores2[i] = int(rand.Int31n(1000)) + 10
	}
	sort.Ints(scores2)
	fmt.Println(scores2[:6])
	worst := make([]int, 5)
	copy(worst, scores2[:6])
	fmt.Println(worst)
}

func after(d time.Duration) chan bool {
	c := make(chan bool)
	go func() {
		time.Sleep(d)
		c <- true
	}()
	return c
}
