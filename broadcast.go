package main

type broadcast struct {
	recvChans []chan []byte
	sendChans []chan []byte
}

func newBroadcast() *broadcast {
	return &broadcast{
		recvChans: make([]chan []byte, 0),
		sendChans: make([]chan []byte, 0),
	}
}

func (b *broadcast) addActor(recvChan chan []byte, sendChan chan []byte) {
	b.recvChans = append(b.recvChans, recvChan)
	b.sendChans = append(b.sendChans, sendChan)
}

func (b *broadcast) run() {
	// Create a channel to synchronize goroutines
	done := make(chan struct{})

	// Launch a goroutine for each sender channel
	for sendIndex, sendChan := range b.sendChans {
		go func(sc <-chan []byte, senderIndex int) {
			defer func() {
				done <- struct{}{}
			}()
			for msg := range sc {
				// Broadcast the message to all receiver channels except the one with the same index
				for recvIndex, recvChan := range b.recvChans {
					if recvIndex != senderIndex {
						// Send the message to the receiver channel
						recvChan <- msg
					}
				}
			}
		}(sendChan, sendIndex)
	}

	// Wait for all sender goroutines to finish
	for range b.sendChans {
		<-done
	}
}
