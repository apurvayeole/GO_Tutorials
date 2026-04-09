package main 
import (
	"fmt"
	"time"
	"sync"
)

func worker (id int, jobs <-chan int,results chan<- int, wg *sync.WaitGroup){
	defer wg.Done() //tells: “this worker has finished” . defer ensures it runs when function exits
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n",id,job)
		time.Sleep(time.Second) //simulate work
		results <- job * 2 //sends result
	}
}

func main(){
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup //Initializes WaitGroup

	//create 3 workers
	for w:=1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs,results, &wg)
	}
	//Now 3 workers are: Running concurrently Waiting for jobs from channel

	for j:= 1; j<=5; j++ {
		jobs<-j
	}
	close(jobs) //Signals: “No more jobs will come”

	//wait for workers
	go func(){
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println("Result:",res)
	}
	wg.Wait() //waits for all workers
}