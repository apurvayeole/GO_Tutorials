package main //entry point package

import ("fmt" //standard library
"time"
)


func add(a int, b int) int{ //function cant be present inside a main function
		return a+b
	}

func divide(a int,b int)(int, int) { //function with multiple return values
	return a/b, a%b
}	

//Structs
type User struct {
	Name string
	Age int
}

//use this function for concurrency
func printMsg(){
	fmt.Println("Hello from goRoutine")
}

func main() {
    fmt.Println("Hello, World!")

	//variables
	var name string = "Apurva"
	age := 20 //infers type automatically

	fmt.Println(name,age)

	//if-else
	if age >= 18 {
		fmt.Println("Adult")
	}else{
		fmt.Println("Minor")
	}

	//loops
	for i:=0; i<5; i++{
		fmt.Println(i)
	}

	result := add(2,4)
	fmt.Println(result)

	fmt.Println(divide(4,2))

	//arrays
	nums := []int{1,2,3}
	nums = append(nums, 4) //adding vlaue at last

	fmt.Println(nums)

	u := User{Name:"Apurva", Age:21}
		fmt.Println(u.Name)

	//concurrency
	go printMsg() //keyword "go" is used to run concurrently
	time.Sleep(time.Second)//without this main() finishes the execution immediately,
	// goroutine may not get time to run
	//time.second give temporary delay of 1 second
}