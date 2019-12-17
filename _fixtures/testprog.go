package main

import (
	"fmt"
	"time"
	"os"
)

func sleepytime() {
	time.Sleep(time.Millisecond)
}

func helloworld() {
	fmt.Println("Hello, World!", "pid:", os.Getpid())
}

func main() {
	for {
		sleepytime()
		helloworld()
	}
}
