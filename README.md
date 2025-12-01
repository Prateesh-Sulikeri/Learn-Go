# Go Concurrency Concepts — Compact Reference

## 📚 Table of Contents & One-Line Definitions

- **WaitGroups** — coordinate and wait for a group of goroutines to finish.
- **Channels**
  - **Unbuffered channels** — operations block until the other side is ready.
  - **Buffered channels** — allow limited sends without waiting.
- **Select** — wait on multiple channel operations.
- **Select with timeout** — avoid blocking forever using `time.After`.
- **Context**
  - **WithCancel** — propagate manual cancellation.
  - **WithTimeout** — auto-cancel after a deadline.
- **Worker Pools** — fixed number of workers processing jobs.
- **Fan-Out / Fan-In** — parallelize work and merge results.
- **Pipelines** — multi-stage concurrent processing.
- **Mutex** — exclusive access to shared state.
- **RWMutex** — multiple readers or one writer.
- **Atomics** — lock-free operations on integers.

---

## WaitGroups
```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        fmt.Println("Hello from goroutine")
        wg.Done()
    }()

    fmt.Println("Hello from main")
    wg.Wait()
}
```

## Unbuffered Channels
```go
package main

import "fmt"

func main() {
    ch := make(chan int)

    go func() {
        fmt.Println("worker: starting")
        ch <- 99
        fmt.Println("worker: done")
    }()

    fmt.Println("main: waiting")
    i := <-ch
    fmt.Println("main: got", i)
}
```

## Buffered Channels
```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    go func() {
        ch <- 10
        ch <- 20
        ch <- 30
    }()
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
```

## Select
```go
package main

import "fmt"

func main() {
    ch1 := make(chan int, 1)
    ch2 := make(chan int, 1)

    ch1 <- 2
    ch2 <- 3

    select {
    case v := <-ch1:
        fmt.Println("ch1", v)
    case v := <-ch2:
        fmt.Println("ch2", v)
    }
}
```

## Select Timeout
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan int)

    go func() {
        time.Sleep(3 * time.Second)
        ch1 <- 1
    }()

    select {
    case v := <-ch1:
        fmt.Println("Value received:", v)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout!")
    }
}
```

## Context WithCancel
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("worker: done!")
                return
            default:
                fmt.Println("worker: doing task!")
                time.Sleep(300 * time.Millisecond)
            }
        }
    }()

    time.Sleep(1 * time.Second)
    cancel()
    time.Sleep(1 * time.Second)
}
```

## Context WithTimeout
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("worker: returning")
                return
            default:
                fmt.Println("worker: working!")
                time.Sleep(300 * time.Millisecond)
            }
        }
    }()

    time.Sleep(2 * time.Second)
}
```

## Worker Pool
```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int)
    results := make(chan int)

    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }

    go func() {
        for i := 1; i <= 5; i++ {
            jobs <- i
        }
        close(jobs)
    }()

    go func() {
        wg.Wait()
        close(results)
    }()

    for r := range results {
        fmt.Println("Result:", r)
    }
}
```

## Fan-Out / Fan-In
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(n int, out chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    time.Sleep(time.Duration(n) * 200 * time.Millisecond)
    out <- n * n
}

func main() {
    out := make(chan int)
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, out, &wg)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    for r := range out {
        fmt.Println("Result:", r)
    }
}
```

## Pipelines
```go
package main

import "fmt"

func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    c := gen(1, 2, 3, 4, 5)
    out := sq(c)

    for r := range out {
        fmt.Println("Result:", r)
    }
}
```
