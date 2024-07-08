package errors

import (
	"fmt"
	"strings"
)

type ErrorChannel chan error

func NewErrorChannel(cap int) ErrorChannel {
	return make(ErrorChannel, cap)
}

func (ec ErrorChannel) Push(err error) {
	if err == nil {
		return
	}
	ec <- err
}

func (ec ErrorChannel) Combined() error {
	if len(ec) == 0 {
		return nil
	}

	strs := make([]string, 0, len(ec))
	for i := 0; i < len(ec); i++ {
		err := <-ec
		if err == nil {
			continue
		}
		strs = append(strs, err.Error())
	}
	return fmt.Errorf(strings.Join(strs, "\n"))
}
