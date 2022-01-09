package errno

import "fmt"

type Err struct {
	Code    int
	Message string
	Err     error
}

func (err *Err) Add(message string) error {
	err.Message += "" + message
	return err.Err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += "" + fmt.Sprintf(format, args...)
	return err.Err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}

// func isErrUserNotFound(err error) bool {
// 	code, _ := DecodeErr(err)
// 	return code == ErrUserNotFound.Code
// }
