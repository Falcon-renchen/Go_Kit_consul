package Services

//第二层
type UserRequest struct {
	Uid int `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}
