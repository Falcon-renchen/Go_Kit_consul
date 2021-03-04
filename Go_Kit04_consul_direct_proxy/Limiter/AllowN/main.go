package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	r := rate.NewLimiter(1,5)
	for {
		if r.AllowN(time.Now(),2) {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("too many request")
		}
		time.Sleep(time.Second)
	}
}
/*
2021-03-03 23:50:26
2021-03-03 23:50:27
2021-03-03 23:50:28
2021-03-03 23:50:29
too many request
2021-03-03 23:50:31
too many request
2021-03-03 23:50:33
*/
