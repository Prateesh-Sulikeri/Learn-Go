// package main

// import (
// 	"fmt"
// )

// func main() {
// 	go func() {
// 		fmt.Println("Hello from goroutine")
// 	}()

// 	fmt.Println("Hello from main")
// 	// time.Sleep(100 * time.Millisecond)
// }

// ===========================================================
// WaitGroups
// ===========================================================

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	func() {
// 		fmt.Println("Hello from goroutine")
// 		wg.Done()
// 	}()

// 	fmt.Println("Hello from main")
// 	wg.Wait()
// }

// ===========================================================
// WaitGroup Exercise
// ===========================================================
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var wg sync.WaitGroup

// 	wg.Add(3)

// 	go func() {
// 		fmt.Println("Goroutine X starting")
// 		fmt.Println("Goroutine X done")
// 		defer wg.Done()
// 	}()

// 	go func() {
// 		fmt.Println("Goroutine X starting")
// 		fmt.Println("Goroutine X done")
// 		defer wg.Done()
// 	}()

// 	go func() {
// 		fmt.Println("Goroutine X starting")
// 		fmt.Println("Goroutine X done")
// 		defer wg.Done()
// 	}()

// 	wg.Wait()
// }

// ===========================================================
// Unbuffered Channels
// ===========================================================

// package main

// import "fmt"

// func main() {
// 	ch := make(chan int)

// 	go func() {
// 		fmt.Println("worker: starting ")
// 		ch <- 99
// 		fmt.Println("worker: done")
// 	}()

// 	fmt.Println("main:waiting")
// 	i := <- ch
// 	fmt.Println("main: got ", i)
// }

// ===========================================================
// Buffered Channels
// ===========================================================

// package main

// import "fmt"

// func main() {

// 	ch := make(chan int, 2)
// 	go func() {
// 		ch <- 10
// 		ch <- 20
// 		ch <- 30
// 	}()

// 	fmt.Println(<-ch)
// 	fmt.Println(<-ch)
// }

// ===========================================================
// Select {}
// ===========================================================

// package main

// import "fmt"

// func main() {
// 	ch1 := make(chan int, 1)
// 	ch2 := make(chan int, 1)

// 	ch1 <- 2
// 	ch2 <- 3

// 	select {
// 		case v := <-ch1:
// 			fmt.Println("ch1", v)
// 		case v := <-ch2:
// 			fmt.Println("ch2", v)
// 	}

// }

// ===========================================================
// Select {} with Timeout
// ===========================================================

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	ch1 := make(chan int)

// 	go func() {
// 		time.Sleep(3 * time.Second)
// 		ch1 <- 1
// 	}()

// 	select {
// 	case v:= <-ch1:
// 		fmt.Println("Value received: ", v)
// 	case <-time.After(1 * time.Second):
// 		fmt.Println("timeout!")
// 	}
// }

// ===========================================================
// Context
// ===========================================================

// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println("worker: done!")
// 				return
// 			default:
// 				fmt.Println("worker: doing task!")
// 				time.Sleep((300 * time.Millisecond))
// 			}
// 		}
// 	}()

// 	time.Sleep(1 * time.Second)
// 	fmt.Println("main: start")
// 	cancel()

// 	time.Sleep(1 * time.Second)
// 	fmt.Println("main: done!")

// }

// ===========================================================
// Context withTimeout
// ===========================================================

// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
// 	defer cancel()

// 	go func() {
// 		for {
// 			select{
// 			case <-ctx.Done():
// 				fmt.Println("worker: returning")
// 				return
// 			default:
// 				fmt.Println("worker: working!")
// 				time.Sleep(300 * time.Millisecond)
// 			}
// 		}
// 	}()

// 	time.Sleep(2 * time.Second)
// }

// ===========================================================
// Worker Pool
// ===========================================================

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func workers(id int, jobs <-chan int, results chan <- int, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for j := range jobs {
// 		results <- j * 2
// 	}
// }

// func main() {
// 	jobs := make(chan int)
// 	results := make(chan int)

// 	var wg sync.WaitGroup

// 	for i:= 1 ; i<=3 ; i++ {
// 		wg.Add(1)
// 		go workers(i, jobs, results, &wg)
// 	}

// 	go func() {
// 		for i:=1 ; i<=5 ; i++ {
// 			jobs <- i
// 		}
// 		close(jobs)
// 	}()

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for r := range results {
// 		fmt.Println("Results: ", r)
// 	}
// }

// ===========================================================
// Fan Out/ Fan In
// ===========================================================

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// ------------------------------------------------------------
// WORKER FUNCTION
//
// Each worker is a FAN-OUT goroutine.
// Many workers run in parallel and send their results into 'out'.
// ------------------------------------------------------------
// func worker(n int, out chan<- int, wg *sync.WaitGroup) {
// 	defer wg.Done()                 // signal this worker is done
// 	time.Sleep(time.Duration(n) * 200 * time.Millisecond)
// 	out <- n * n                    // send result into merged channel
// }

// func main() {

// 	out := make(chan int)           // MERGED results channel

// 	var wg sync.WaitGroup           // Wait for workers to finish

// 	// ------------------------------------------------------------
// 	// FAN-OUT: start 3 workers in parallel
// 	//
// 	// All workers write to the SAME channel: out
// 	//
// 	// This is the FAN-OUT part (spawning parallel workers).
// 	// ------------------------------------------------------------
// 	for i := 1; i <= 3; i++ {
// 		wg.Add(1)                    // tell the WaitGroup a worker is starting
// 		go worker(i, out, &wg)       // start worker i
// 	}

// 	// ------------------------------------------------------------
// 	// CLOSE RESULTS ONLY AFTER ALL WORKERS FINISH
// 	//
// 	// This goroutine is the "closing coordinator".
// 	// It waits for all workers (wg.Wait)
// 	// and then closes the out channel.
// 	//
// 	// Once out is closed, the FAN-IN loop stops.
// 	// ------------------------------------------------------------
// 	go func() {
// 		wg.Wait()
// 		close(out)
// 	}()

// 	// ------------------------------------------------------------
// 	// FAN-IN: collect results from ALL workers
// 	//
// 	// This loop reads from 'out' channel.
// 	// It merges values from all workers.
// 	// When the channel closes, the loop ends.
// 	//
// 	// THIS is the FAN-IN portion.
// 	// ------------------------------------------------------------
// 	for r := range out {
// 		fmt.Println("Result:", r)
// 	}
// }

// ===========================================================
// Pipelines
// ===========================================================

package main

import "fmt"

func gen(nums ... int) <- chan int {
	out := make(chan int)

	go func ()  {
		for _, n := range nums {
			out <- n 
		}	

		close(out)
	}()
	return out
}

func double(in <- chan int) <- chan int {
	out := make(chan int)
	
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func sq(in <- chan int) <- chan int {
	out := make(chan int)
	go func ()  {
		for n := range in {
			out <- n * n 
		}	
		close(out)
	}()
	return out
}

func main() {
	c := gen(1, 2, 3, 4, 5)
	
	dbl := double(c)
	
	out := sq(dbl)

	// Print the contents 
	for r := range(out) {
		fmt.Println("Result: ", r)
	}
}