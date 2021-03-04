package main

import (
	"golang.org/x/time/rate"
	"net/http"
)

//限流demo

var r = rate.NewLimiter(1,5)	//   5个token， 每次拿出1个   空了的话输出too many request

func MyLimit(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !r.Allow() {
			http.Error(writer, "too many request", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer,request)
	})
}


func main() {


	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK!!!"))
	})


	http.ListenAndServe(":8080",MyLimit(mux))
}
