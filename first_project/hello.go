package main

import "fmt"

func producer(ch chan int) {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch)
}


func main() {
    // ch := make(chan int) //channel that sends/receive values

    // go func() { //goroutine
    //     ch <- 42   // send value
    // }()

    // value := <-ch  // receive value
	// //data flows in the direction of arrow
    // fmt.Println(value)

	ch2 := make(chan int)
	go producer(ch2)

	for val := range ch2 {
		fmt.Println(val)
	}

}