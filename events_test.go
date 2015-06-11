package main

import (
	"github.com/cheekybits/is"
	"testing"
	"time"
)

func TestWhen(t *testing.T) {

	is := is.New(t)
	actionCalls := 0
	condition := false
	trigger := make(chan interface{})

	triggerFunc := func(events chan<- interface{}) {
		for range trigger {
			events <- true
		}
	}
	conditionFunc := func() bool {
		return condition
	}
	actionFunc := func() {
		actionCalls += 1
	}

	When(triggerFunc, conditionFunc, actionFunc)
	is.Equal(actionCalls, 0)

	trigger <- true
	time.Sleep(10 * time.Millisecond)
	is.Equal(actionCalls, 0)

	condition = true
	trigger <- true
	time.Sleep(10 * time.Millisecond)
	is.Equal(actionCalls, 1)

}
