package Services

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
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

func GenUserEndpoint(userSerivce IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := userSerivce.GetName(r.Uid)
		if r.Method=="GET" {
			result = userSerivce.GetName(r.Uid)
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
