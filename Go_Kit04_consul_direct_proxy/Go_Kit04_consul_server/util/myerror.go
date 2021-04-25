package util

//自定义error
type MyError struct {
	Code int
	Message string
}

func (this *MyError) Error() string {
	return this.Message
}

func NewMyError(code int, message string) error {
	return &MyError{
		Code:    code,
		Message: message,
	}
}

