package main

import (
	"fmt"
	"github.com/koding/kite"
	"github.com/koding/kite/config"
	"net/url"
)
func main() {
	//创建一个新的kite，参数为名称和服务版本
	k := kite.New("math","1.0.0")
	c := config.MustGet()
	k.Config = c

	//为kontrol服务器设置正确的URI
	k.Config.KontrolURL = "http://kontrol:6000/kite"
	k.RegisterForever(&url.URL{
		//http是主机名
		Scheme:     "http",
		Opaque:     "127.0.0.1",
		User:       nil,
		Host:       "",
		Path:       "/kite",
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   "",
		Fragment:   "",
	})
	//设置路由，传递一个负责执行该请求的参数。Args是一个dnode消息，
	k.HandleFunc("Hello", func(r *kite.Request) (interface{},error) {
		name, _ := r.Args.One().String()
		return fmt.Sprintf("Hello %v", name), nil
	}).DisableAuthentication()

	//放入8091端口，并且run
	k.Config.Port = 8091
	k.Run()
}


//该Dial()方法需要超时，超时后，无论是否有可能连接到下游服务，该方法都会返回。该DialForever()方法，作为方法顾名思义，将不会返回。
//在这两种情况下，都会返回一个通道，我们将使用该通道暂停执行直到获得连接。
//
//现在，调用该服务非常简单Tell，只需执行，传递要运行的方法的名称以及将参数作为该方法的接口传递。以我的拙见，风筝在这里丢分。
//服务呼叫的合同非常宽松，为消费者创建实现将是不费吹灰之力的。
