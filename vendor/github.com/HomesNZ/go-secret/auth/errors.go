package auth

import (
	"bytes"
	"errors"
)

var (
	errAuthInvalid = errors.New("auth token invalid")
)

// multiError is a simple way to collect multiple errors together to return as
// a single error. Each message is separated/suffixed with "; ".
type multiError []error

func (es multiError) Error() string {
	b := bytes.Buffer{}
	for _, e := range es {
		b.WriteString(e.Error())
		b.WriteString("; ")
	}
	return b.String()
}
