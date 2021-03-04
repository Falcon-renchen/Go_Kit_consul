package main

import (
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_client_direct/Services"
	"context"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/url"
)

//直连服务端

func main() {
	//先写死
	tgt, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	//	此api专门给客户端调用
	//创建一个直连client，如何请求，如何回复
	client := httptransport.NewClient("GET",tgt, Services.GetUserInfo_Request, Services.GetUserInfo_Response)
	//暴露出endpoint，然后执行
	getUserInfo := client.Endpoint()

	ctx := context.Background()//传一个空的上下文对象
	//执行
	res, err := getUserInfo(ctx, Services.UserRequest{
		Uid:    101,
		Method: "",
	})
	if err != nil {
		panic(err)
	}
	//因为res是interface，所以断言
	userInfo := res.(Services.UserResponse)
	fmt.Println(userInfo.Result)

}