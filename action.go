package main
import (
	"log"
)

type Action interface {
	When(event <-chan bool)
}

type LogAction struct{}

func (l *LogAction)When(event <-chan bool) {
	for range event {
		log.Println("Hello")
	}
}
