package main

import (
	"fmt"
	"time"
)

func main() {
	tm := time.Now()
	fmt.Println("Hello, 世界", tm.UTC())
}
