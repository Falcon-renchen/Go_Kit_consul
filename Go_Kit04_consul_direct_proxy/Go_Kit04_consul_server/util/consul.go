package util

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	"log"
)


var ConsulClient *consulapi.Client
var ServiceID string
var ServiceName string
var ServicePort int

func init() { //引入包时自动导入
	config := consulapi.DefaultConfig()
	config.Address = "172.16.17.151:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceID = "userservice" + uuid.New()
}

func SetServiceAndPort(name string, port int)  {
	ServiceName = name
	ServicePort = port
}

//注册consul
func RegService()  {
	//config := consulapi.DefaultConfig()
	//
	//config.Address = "172.16.17.150:8500"
	reg := consulapi.AgentServiceRegistration{
		Kind:              "",
		ID:                ServiceID,
		Name:              ServiceName,
		Tags:              []string{"primary"},
		Port:              ServicePort,
		Address:           "192.168.0.107",
		TaggedAddresses:   nil,
		EnableTagOverride: false,
		Meta:              nil,
		Weights:           nil,
		Check:             nil,
		Checks:            nil,
		Proxy:             nil,
		Connect:           nil,
	}
	check := consulapi.AgentServiceCheck{
		CheckID:                        "",
		Name:                           "",
		Args:                           nil,
		DockerContainerID:              "",
		Shell:                          "",
		Interval:                       "5s",
		Timeout:                        "",
		TTL:                            "",
		Header:                         nil,
		Method:                         "",
		TCP:                            "",
		Status:                         "",
		Notes:                          "",
		TLSSkipVerify:                  false,
		GRPC:                           "",
		GRPCUseTLS:                     false,
		AliasNode:                      "",
		AliasService:                   "",
		DeregisterCriticalServiceAfter: "",
	}
	check.HTTP = fmt.Sprintf("http://%s:%d/health",reg.Address,ServicePort)//动态获取ip和port
	reg.Check = &check

	
	err := ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

//反注册
func UnRegService()  {
	ConsulClient.Agent().ServiceDeregister(ServiceID)
}