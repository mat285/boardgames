package errors

import "sync"

type ErrorWaitGroup struct {
	ErrorChannel
	sync.WaitGroup
}

func NewErrorWaitGroup(cap int) *ErrorWaitGroup {
	return &ErrorWaitGroup{
		ErrorChannel: NewErrorChannel(cap),
	}
}

func (wg *ErrorWaitGroup) PushDone(err error) {
	wg.ErrorChannel.Push(err)
	wg.Done()
}

func (wg *ErrorWaitGroup) Wait() error {
	wg.WaitGroup.Wait()
	close(wg.ErrorChannel)
	return wg.Combined()
}
