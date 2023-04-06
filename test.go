package main

import (
	"fmt"
	"sync"
)

func test() {
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()

			Print(i)
			//fmt.Println(modules[i])
			//
			//if modules[i] == nil {
			//	return
			//}
			//
			//modules[i].Detect( item)
		}()
	}

	wg.Wait()
}

func Print(i int) {
	fmt.Println(i)
}
