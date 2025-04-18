package main

import (
	"fmt"
	"runtime"
	"sync"
)

func WaitGroupDemo() {
	numbers := []int{10, 20, 30, 40, 50}
	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)
		fmt.Println(runtime.NumGoroutine())
		go func() {
			num += 1
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("main function")
	fmt.Println(runtime.NumGoroutine())
}
