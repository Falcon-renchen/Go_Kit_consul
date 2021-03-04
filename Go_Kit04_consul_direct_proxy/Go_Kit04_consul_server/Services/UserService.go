package Services

import "errors"

//第一层
type IUserService interface {
	GetName(userid int) string
	DelUser(userid int) error
}

type UserService struct {

}

func (this UserService) GetName(userid int) string {
	if userid == 101 {
		return "shenyi"
	}
	return "guest"
}


//假设性的管理员
func (this UserService) DelUser(userid int) error {
	if userid == 101 {
		return errors.New("无权限")
	}
	return nil
}
