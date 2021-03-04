package Services

import (
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_server/util"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"strconv"
)
import "context"

//第二层
type UserRequest struct {
	Uid int `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}

//日志中间件
func UserServiceLogMiddleware(logger kitlog.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method",r.Method,"event","get user","userid",r.Uid)//日志输出
			return next(ctx,request)
		}
	}
}

//加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("too many errors")
			}
			return next(ctx,request)
		}
	}
}

func GenUserEndpoint(userSerivce IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		r := request.(UserRequest)
		result := userSerivce.GetName(r.Uid)
		if r.Method=="GET" {
			result = userSerivce.GetName(r.Uid)+strconv.Itoa(util.ServicePort)
			//logger.Log("method",r.Method,"event","get user","userid",r.Uid)//日志输出
		} else if r.Method == "DELETE" {//代表删除
			err := userSerivce.DelUser(r.Uid)
			if err != nil {//代表有错
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户删除成功",r.Uid)
			}
		}
		return UserResponse{Result:result}, nil
	}
}
