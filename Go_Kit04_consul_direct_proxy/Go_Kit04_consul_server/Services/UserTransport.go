package Services

import (
	"Go_kit/Go_Kit04_consul_direct_proxy/Go_Kit04_consul_server/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	//http://localhost:xxx/?uid=101
	//if r.URL.Query().Get("uid")!="" {
	//	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	//	return UserRequest{
	//		Uid:uid,
	//	}, nil
	//}
	vars := mux.Vars(r)
	if uid, ok := vars["uid"];ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			Uid:uid,
			Method: r.Method,
		},nil
	}
	

	return nil, errors.New("参数错误")
}

func EncodeUserResponse(c context.Context,w http.ResponseWriter, response interface{})  error {
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter)  {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type",contentType)
	if myerr, ok := err.(*util.MyError);ok {
		w.WriteHeader(myerr.Code)
		w.Write(body)
	} else {
		w.WriteHeader(500)
		w.Write(body)
	}


}