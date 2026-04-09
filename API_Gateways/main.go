//Stage 1

// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func handler(w http.ResponseWriter, r *http.Request){
// 	fmt.Fprintln(w, "Hello from go server")
// }

// func main(){
// 	http.HandleFunc("/",handler) // request st / called handler function
// 	http.ListenAndServe(":8000",nil) //starts the server
// }

package main 

import(
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// func handler(w http.ResponseWriter, r *http.Request){
 	// fmt.Fprintln(w, "Hello from go server")
// }
func main(){
	//target service
	target, _ := url.Parse("http://localhost:3001") //request apperars here

	//create proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request){
		log.Println("Forwarding to user service")
		proxy.ServeHTTP(w,r)
	})

	log.Println("gateway running on 3000")
	http.ListenAndServe(":3000",nil)
}
