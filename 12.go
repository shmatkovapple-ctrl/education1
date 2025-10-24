package main

import "fmt"

// определите функцию download()
func download(s1, s2, s3 int, ch1, ch2, ch3 chan int) {
	var result1, result2, result3 int

	for i := 0; i <= s1; i++ {
		result1 += i
	}
	for j := 0; j <= s2; j++ {
		result2 += j
	}
	for k := 0; k <= s3; k++ {
		result3 += k
	}
	ch1 <- result1
	ch2 <- result2
	ch3 <- result3
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	var s1, s2, s3 int
	fmt.Scanln(&s1)
	fmt.Scanln(&s2)
	fmt.Scanln(&s3)

	go download(s1, s2, s3, ch1, ch2, ch3)

	//выведите сумму всех результатов
	fmt.Println(<-ch1 + <-ch2 + <-ch3)
}
