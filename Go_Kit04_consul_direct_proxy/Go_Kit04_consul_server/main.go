package main

import (
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_server/Services"
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_server/util"
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// go run main.go --name userservice -p 8080/8081
	name := flag.String("name","","服务名称")
	port := flag.Int("p",0,"服务端口")
	flag.Parse()
	if *name=="" {
		log.Fatal("请指定服务名称")
	}
	if *port==0 {
		log.Fatal("请指定端口")
	}
	util.SetServiceAndPort(*name,*port) //动态设置服务名和端口


	user := Services.UserService{}
	limit := rate.NewLimiter(1,3)

	//日志输出
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.With(logger,"time",kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger,"caller",kitlog.DefaultCaller)
	}

	endp := Services.RateLimit(limit)(Services.UserServiceLogMiddleware(logger)(Services.GenUserEndpoint(user))) //先过ratelimit

	//自定义错误处理
	options := []httptransport.ServerOption {
		httptransport.ServerErrorEncoder(Services.MyErrorEncoder),
	}

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse,options...)

	//路由
	router := mux.NewRouter()
	//r.Handle(`user/{uid\d+}`,serverHandler)
	{
		router.Methods("GET","DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		//不使用三层架构，直接写死一个路由
		router.Methods("GET").Path(`/health`).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type","application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})

	}
	errChannel := make(chan error)

	go func() {
		util.RegService() //注册到consul
		//仅支持get的方法请求
		err := http.ListenAndServe(":"+strconv.Itoa(*port),router)
		if err != nil {
			log.Fatal(err)
			errChannel <- err
		}
	}()
	go func() {
		sigChannel := make(chan os.Signal)
		signal.Notify(sigChannel,syscall.SIGINT,syscall.SIGTERM)
		errChannel <- fmt.Errorf("%s",<-sigChannel)
	}()

	getErr := <-errChannel
	util.UnRegService()
	log.Fatal(getErr)

}
