package main

/*
func chanEcho(ch chan any, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- "juniper"
	ch <- []int{10, 20, 30}
	ch <- time.Now()
}
*/

/*
func main() {

		var wg sync.WaitGroup

		ch := make(chan any, 15)
		wg.Add(3)

		for range 3 {
			go chanEcho(ch, &wg)
		}
		wg.Wait()
		fmt.Println(len(ch))

}
*/

import (
	"log"
)

func main() {
	ch := make(chan int, 1)
	go dummy(ch)
	log.Println("waiting for reception...")
	ch <- 45
	ch <- 58
	ch <- 100
}

func dummy(c chan int) {
	smth := <-c
	log.Println("has received something", smth)
}
