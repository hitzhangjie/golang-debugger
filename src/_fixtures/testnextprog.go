package main

import (
	"fmt"
	"os"
	"time"
)

func sleepytime() {
	time.Sleep(time.Millisecond)
}

func helloworld() {
	fmt.Println("Hello, World! pid:", os.Getpid())
}

func testnext() {
	j := 1

	for i := 0; i <= 1; i++ {
		j += j * (j ^ 3) / 100

		helloworld()
	}


}

func main() {
	for {
		sleepytime()
		testnext()
	}
}
