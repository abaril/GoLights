package main


type TriggerFunc func(chan<- interface{})

type ConditionFunc func() bool

type ActionFunc func()


func When(trigger TriggerFunc, condition ConditionFunc, action ActionFunc) {

	events := make(chan interface{})
	go trigger(events)
	go func() {
		for {
			select {
			case <-events:
				if condition() {
					action()
				}
			}
		}
	}()
}

