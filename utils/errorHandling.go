package utils

import "fmt"

func HandleError(err error) {
	fmt.Println(err.Error())
}

type ErrATagNotFound struct {
}

func (e *ErrATagNotFound) Error() string {
	return fmt.Sprintf("not a tag")
}
