package main

import ap "github.com/LiangrunDa/AutomergePrototype"

type Checker interface {
	getInitialData() (*ap.Automerge, InitialMetadata)
	initialize(doc *ap.Automerge, change []byte, meta InitialMetadata)
	handleMessage(doc *ap.Automerge, message []byte)
	do(doc *ap.Automerge)
	getNumberOfOperations() int
	getMaxInterval() int
}

type InitialMetadata = map[string]any

type Parameters struct {
	NumberOfOperations int
	MaxInterval        int
}

func NewParameters(numberOfOperations int, maxInterval int) *Parameters {
	return &Parameters{
		NumberOfOperations: numberOfOperations,
		MaxInterval:        maxInterval,
	}
}

func (p *Parameters) getNumberOfOperations() int {
	return p.NumberOfOperations
}

func (p *Parameters) getMaxInterval() int {
	return p.MaxInterval
}
