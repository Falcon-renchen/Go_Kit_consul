package main

import (
	"Go_kit/Go_Kit03_struct/Services"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndpoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse)

	//路由
	router := mux.NewRouter()
	//r.Handle(`user/{uid\d+}`,serverHandler)
	{
		router.Methods("GET","DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		//不使用三层架构，直接写死一个路由
		router.Methods("GET").Path(`/health`).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type","application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})

	}
	//仅支持get的方法请求
	http.ListenAndServe(":8080",router)

}
