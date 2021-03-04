package account

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	CreateUserRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	CreateUserResponse struct {
		ok string `json:"ok"`
	}
	GetUserResponse struct {
		Email string `json:"email"`
	}
	GetUserRequset struct {
		Id string `json:"id"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeUserReq(ctx context.Context,r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeEmailReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRequset

	vars := mux.Vars(r)

	req = GetUserRequset{Id:vars["id"]}
	return req, nil
}