package util

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	"net/url"
	"os"
	"time"
	"fmt"
	"context"
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_client_proxy/Services"
)
// string : 从服务端获取的结果是什么
// error  : 返回错误
func GetUser() (string, error) {
	//第一步创建client
	config := consulapi.DefaultConfig()
	config.Address = "172.16.17.151:8500"//注册中心地址,,虚拟机的地址
	api_client, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	client := consul.NewClient(api_client)


	//查看服务实例状态
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	tags := []string{"primary"}
	instancer := consul.NewInstancer(client, logger, "userservice", tags, true)


	factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {
		tart, _ := url.Parse("http://"+service_url) //  192.168.0.107:8080  真实服务的地址
		return httptransport.NewClient("GET", tart, Services.GetUserInfo_Request, Services.GetUserInfo_Response).Endpoint(),nil,nil
	}

	//获取所有端点
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	endpoints, err := endpointer.Endpoints()
	fmt.Println("服务有",len(endpoints),"条")


	//轮询算法,负载均衡
	//mylb := lb.NewRoundRobin(endpointer)

	//随机算法
	mylb := lb.NewRandom(endpointer, time.Now().UnixNano())

	for {
		//getUserInfo := endpoints[0] 	//	写死获得第一个服务
		getUserInfo, _ := mylb.Endpoint()   //根据算法获取服务


		ctx := context.Background()//传一个空的上下文对象
		//执行
		res, err := getUserInfo(ctx, Services.UserRequest{
			Uid:    101,
			Method: "",
		})
		if err != nil {
			panic(err)
			os.Exit(1)
		}
		//因为res是interface，所以断言
		userInfo := res.(Services.UserResponse)
		fmt.Println(userInfo.Result)

		time.Sleep(time.Second*3)
	}

}
