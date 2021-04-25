package main

import (
	"Go_kit/Go_Kit02/account"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const dbsource = "postgresql://postgres:postgres@localhost:5432/gokitexample?sslmode=disable"


func main() {
	var httpAddr = flag.String("http",":8080","http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"server", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	level.Info(logger).Log("msg","service started")
	defer level.Info(logger).Log("msg","service ended")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("postgres", dbsource)
		if err != nil {
			level.Error(logger).Log("exit",err)
			os.Exit(-1)
		}
	}
	flag.Parse()
	ctx := context.Background()	//构建为服务
	var srv account.Service
	{
		//创建一个库，用来调用数据库
		repository := account.NewRepo(db, logger)

		srv = account.NewService(repository, logger)
	}
	//建立一个错误通道，终止的时候会有信号，记录错误
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	//将服务传递给端点
	endpoints := account.MakeEndpoints(srv)
	//调用新的http
	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}