# bchan
Better Channels for Golang

## Example Usage

```go
package main

import (
    "fmt"
    "github.com/B00TK1D/bchan"
)

func main() {
    in, out := BChan[int]{}.Unbounded()
	in <- 1
	in <- 2

	fmt.Println(<-out)
	fmt.Println(<-out)

	in <- 3
	fmt.Println(<-out)
}
```

## Performance


This library enables both lazily allocated and unbounded channels, leading to significant memory savings for use cases involving many channels for which the channel size is unknown at compile time.
Benchmarks indicate roughly 6x CPU performance loss compared to native channels for both the bounded and unbounded channel implementations.
