package Services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

//客户端的server
//和server端的相反

func GetUserInfo_Request(ctx context.Context,request *http.Request,r interface{} ) error {
	user_request := r.(UserRequest)
	request.URL.Path = "/user/" + strconv.Itoa(user_request.Uid)
	return nil
}

func GetUserInfo_Response(ctx context.Context, res *http.Response) (response interface{}, err error) {
	if res.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var user_response UserResponse
	err = json.NewDecoder(res.Body).Decode(&user_response)//解码res.body的东西然后解码到user_response里面
	if err != nil {
		return nil, err
	}
	return user_response, nil
}