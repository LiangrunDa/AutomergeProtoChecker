package main

import (
	ap "github.com/LiangrunDa/AutomergePrototype"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
)

// Check1 generate 100 objects, then randomly move them
type Check1 struct {
	objects []ap.ExOpId
	*Parameters
}

func NewCheck1() *Check1 {
	return &Check1{Parameters: NewParameters(1000, 100)}
}

func (c *Check1) getInitialData() (*ap.Automerge, InitialMetadata) {
	doc := ap.NewAutomerge(uuid.New())
	meta := InitialMetadata{}
	objects := make([]ap.ExOpId, 100)
	t1 := doc.StartTransaction()
	// generate 100 objects
	for i := 0; i < 100; i++ {
		if i == 0 {
			objects[i], _ = t1.PutObject(ap.ExRootOpId, strconv.Itoa(i), ap.MAP)
		} else {
			r := rand.Intn(i)
			objects[i], _ = t1.PutObject(objects[r], strconv.Itoa(i), ap.MAP)
		}
	}
	doc.CommitTransaction()
	meta["objects"] = objects
	c.objects = objects
	return doc, meta
}

func (c *Check1) initialize(doc *ap.Automerge, change []byte, meta InitialMetadata) {
	c.objects = meta["objects"].([]ap.ExOpId)
	doc.MergeFromChangeBytes(change)
}

func (c *Check1) handleMessage(doc *ap.Automerge, message []byte) {
	doc.MergeFromChangeBytes(message)
}

func (c *Check1) do(doc *ap.Automerge) {
	// generate random move
	t1 := doc.StartTransaction()
	src := rand.Intn(100)
	dst := rand.Intn(100)
	for src == dst {
		dst = rand.Intn(100)
	}
	t1.MoveObject(c.objects[src], c.objects[dst])
	doc.CommitTransaction()
}

// Check2 randomly generate objects
type Check2 struct {
	objects []ap.ExOpId
	*Parameters
}

func NewCheck2() *Check2 {
	return &Check2{Parameters: NewParameters(200, 100)}
}

func (c *Check2) getInitialData() (*ap.Automerge, InitialMetadata) {
	doc := ap.NewAutomerge(uuid.New())
	meta := InitialMetadata{}
	return doc, meta
}

func (c *Check2) initialize(doc *ap.Automerge, change []byte, meta InitialMetadata) {
	// do nothing
}

func (c *Check2) handleMessage(doc *ap.Automerge, message []byte) {
	objects := doc.MergeFromChangeBytesAndGetNewObjects(message)
	c.objects = append(c.objects, objects...)
}

func (c *Check2) do(doc *ap.Automerge) {
	// random generate objects
	t1 := doc.StartTransaction()
	num := len(c.objects)
	if num == 0 {
		t1.PutObject(ap.ExRootOpId, strconv.Itoa(num), ap.MAP)
	} else {
		r := rand.Intn(num)
		object, _ := t1.PutObject(c.objects[r], strconv.Itoa(num), ap.MAP)
		c.objects = append(c.objects, object)
	}

	doc.CommitTransaction()
}

// Check3 randomly generate objects and move them
type Check3 struct {
	objects []ap.ExOpId
	*Parameters
}

func NewCheck3() *Check3 {
	return &Check3{Parameters: NewParameters(200, 100)}
}

func (c *Check3) getInitialData() (*ap.Automerge, InitialMetadata) {
	doc := ap.NewAutomerge(uuid.New())
	meta := InitialMetadata{}
	return doc, meta
}

func (c *Check3) initialize(doc *ap.Automerge, change []byte, meta InitialMetadata) {
	// do nothing
}

func (c *Check3) handleMessage(doc *ap.Automerge, message []byte) {
	objects := doc.MergeFromChangeBytesAndGetNewObjects(message)
	c.objects = append(c.objects, objects...)
}

func (c *Check3) do(doc *ap.Automerge) {
	t1 := doc.StartTransaction()
	op := rand.Intn(2)
	num := len(c.objects)
	randomPropertyName := strconv.Itoa(rand.Intn(100))
	if num < 2 || op == 0 {
		if num == 0 {
			object, _ := t1.PutObject(ap.ExRootOpId, randomPropertyName, ap.MAP)
			c.objects = append(c.objects, object)
		} else {
			r := rand.Intn(num)
			object, _ := t1.PutObject(c.objects[r], randomPropertyName, ap.MAP)
			c.objects = append(c.objects, object)
		}
	} else {
		moved := false
		times := 0
		for moved {
			src := rand.Intn(num)
			dst := rand.Intn(num)
			for src == dst {
				dst = rand.Intn(num)
			}
			err := t1.MoveObject(c.objects[src], c.objects[dst])
			if err == nil {
				moved = true
			}
			times++
			if times > 10 {
				moved = true
			}
		}
	}
	doc.CommitTransaction()
}
