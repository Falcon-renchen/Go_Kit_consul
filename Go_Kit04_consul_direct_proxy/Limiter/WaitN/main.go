package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"time"
)
//限流器
func main() {
	r := rate.NewLimiter(1,5)
	ctx := context.Background()

	for {
		err := r.WaitN(ctx, 2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second)
	}
}
/*
2021-03-03 23:51:04
2021-03-03 23:51:05
2021-03-03 23:51:06
2021-03-03 23:51:07
2021-03-03 23:51:09
2021-03-03 23:51:11
2021-03-03 23:51:13
 */
