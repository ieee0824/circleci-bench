package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var (
	dx = 9
	dy = 3
	dz = 1
)

var writeStr string

func getParam(n int) (int, int, int) {
	return dx * n, dy * n, dz * n
}

func io(n int) {
	file, err := os.OpenFile("test/"+strconv.Itoa(n), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	file.WriteString(writeStr)
	file.Close()
	fmt.Println("write")
}

func ioWorker(wg *sync.WaitGroup, q chan int) {
	defer wg.Done()
	for {
		i, ok := <-q
		if !ok {
			return
		}
		io(i)
	}
}

func tarai(x, y, z int) int {
	if x <= y {
		return y
	}
	return tarai(
		tarai(x-1, y, z),
		tarai(y-1, z, x),
		tarai(z-1, x, y),
	)
}

func worker(wg *sync.WaitGroup, t chan bool) {
	defer wg.Done()

	for {
		_, ok := <-t
		if !ok {
			return
		}
		fmt.Println(tarai(getParam(2)))
	}

	for i := 0; i < 65536; i++ {
		writeStr += "abcedfghijklmnopqrstuvwzyz\n"
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	t := make(chan bool, 2048)
	q := make(chan int, runtime.NumCPU())
	killer := make(chan bool)
	var wg sync.WaitGroup
	var ioW sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(&wg, t)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		ioW.Add(1)
		go ioWorker(&ioW, q)
	}

	go func() {
	loop:
		for {
			select {
			case <-killer:
				close(q)
				break loop

			default:
				for i := 0; i < runtime.NumCPU(); i++ {
					q <- i
				}
			}
		}
	}()

	for i := 0; i < runtime.NumCPU()*2; i++ {
		t <- true
	}

	close(t)
	wg.Wait()
	fmt.Println("end tarai")
	killer <- true
	ioW.Wait()
}
