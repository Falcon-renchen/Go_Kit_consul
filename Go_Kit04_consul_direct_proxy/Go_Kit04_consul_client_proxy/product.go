package main

import (
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_client_proxy/Services"
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_client_proxy/util"
	"io"
	"net/url"
	"os"
	"time"
)

//通过连接consul获取服务
func main_2() {
	//第一步创建client
	config := consulapi.DefaultConfig()
	config.Address = "172.16.17.151:8500"//注册中心地址
	api_client,_ := consulapi.NewClient(config)
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
	endpoints, _ := endpointer.Endpoints()
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


func main()  {
	config := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,//支持最大并发数为5
		RequestVolumeThreshold: 20,//有20个请求才进行错误百分比计算
		SleepWindow:            5,//熔断器：默认关闭，请求次数异常超过设定比例则打开，打开后直接执行降级函数，半开为定期打开
		//过了5s则尝试服务是否可用。默认5s
		ErrorPercentThreshold:  20,//超过20%熔断器打开，错误百分比，默认为50%
	}
	hystrix.ConfigureCommand("getuser",config)
	err := hystrix.Do("getuser", func() error {
		res, err := util.GetUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
		return err
	}, func(err error) error {//有问题就会降级
		fmt.Println("降级用户")
		return err
	})

	if err != nil {
		log.Fatal(err)
	}


}