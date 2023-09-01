package main

import (
	"fmt"
	ap "github.com/LiangrunDa/AutomergePrototype"
	"github.com/google/uuid"
	"sync"
	"time"
)

func main() {
	ap.DebugMode()
	ap.SetLogPath("./test.log")
	// create a broadcast
	b := newBroadcast()
	// create actors
	actors := make([]*Actor, 0)
	for i := 0; i < 3; i++ {
		recvChan := make(chan []byte)
		sendChan := make(chan []byte)
		actor := NewActor(recvChan, sendChan, NewCheck3())
		actors = append(actors, actor)
		b.addActor(recvChan, sendChan)
	}

	// init actors
	initialize(actors)

	// start broadcast
	go b.run()

	var wg sync.WaitGroup

	// start listening
	for _, actor := range actors {
		go actor.listen()
	}

	for _, actor := range actors {
		wg.Add(1)
		actor := actor
		go func() {
			actor.run()
			wg.Done()
		}()
	}

	wg.Wait()                   // Wait for all actors to finish
	time.Sleep(time.Second * 5) // wait for all messages to be broadcasted and processed

	// check if all actors reach the same state
	if ok, r := finalStateCheck(actors); ok {
		fmt.Println("😁SUCCESS: All actors reach the same state")
	} else {
		fmt.Println("😭FAIL: " + r)
	}
}

func initialize(actors []*Actor) {
	doc, meta := actors[0].checker.getInitialData()
	actors[0].doc = doc
	for i := 1; i < len(actors); i++ {
		actors[i].doc = ap.NewAutomerge(uuid.New())
	}

	change, err := doc.GetLatestChangeBytes()
	if err == nil {
		for i := 1; i < len(actors); i++ {
			actors[i].checker.initialize(actors[i].doc, change, meta)
		}
	}
}

func finalStateCheck(actors []*Actor) (bool, string) {
	firstTree := actors[0].doc.GetDocumentTree()
	leftTrees := make([]map[ap.ExOpId]ap.ExOpId, 0)
	for i := 1; i < len(actors); i++ {
		tree := actors[i].doc.GetDocumentTree()
		leftTrees = append(leftTrees, tree)
	}
	for k, v := range firstTree {
		for _, tree := range leftTrees {
			if tree[k] != v {
				return false, fmt.Sprintf("key %v has different value %v, %v", k, v, tree[k])
			}
		}
	}
	return true, ""
}
