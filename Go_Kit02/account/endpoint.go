package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)


//端点
type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:  makeCreateUserEndpoints(s),
		GetUser:     makeGetUserEndpoint(s),
	}
}

func makeCreateUserEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.Email, req.Password)
		return CreateUserResponse{ok:ok}, err
	}
}
//将接口转换成获取用户请求，然后调用服务，让用户获得请求id
func makeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequset)
		email, err := s.GetUser(ctx, req.Id)

		return GetUserResponse{Email:email}, err
	}
}