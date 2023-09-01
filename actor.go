package main

import (
	"fmt"
	ap "github.com/LiangrunDa/AutomergePrototype"
	"math/rand"
	"time"
)

type Actor struct {
	doc     *ap.Automerge
	checker Checker
	recvMsg chan []byte
	sendMsg chan []byte
}

func NewActor(recvMsg chan []byte, sendMsg chan []byte, checker Checker) *Actor {
	return &Actor{
		recvMsg: recvMsg,
		sendMsg: sendMsg,
		checker: checker,
	}
}

func (a *Actor) listen() {
	for msg := range a.recvMsg {
		a.checker.handleMessage(a.doc, msg)
	}
}

func (a *Actor) run() {
	fmt.Print("Actor running\n")
	for i := 0; i < a.checker.getNumberOfOperations(); i++ {
		a.checker.do(a.doc)
		change, _ := a.doc.GetLatestChangeBytes()
		a.sendMsg <- change
		// sleep at most MaxInterval ms
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(a.checker.getMaxInterval())))
	}
	close(a.sendMsg)
}
