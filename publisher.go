package main

import (
	"log"
	"sync"
	"time"

	"github.com/nsqio/nsq"
)

type Publisher interface {
	Publish(Message) error
	Stop() error
}

// NewBufferedNSQPublisher is responsible for buffering messages and flushing
// them in batches to NSQ. This publisher accepts an NSQ producer connection
// and uses that to write messages in batches.
func NewBufferedNSQPublisher(bufferSize int, flushInterval time.Duration, topicName string, producer *nsq.Producer) Publisher {
	publisher := &bufferedNSQPublisher{
		messageCh:     make(chan Message, bufferSize/2),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,

		topicName: topicName,
		producer:  producer,
	}
	publisher.startLoop()
	return publisher
}

type bufferedNSQPublisher struct {
	messageCh     chan Message
	flushInterval time.Duration
	bufferSize    int

	topicName string
	producer  *nsq.Producer
	wg        sync.WaitGroup
}

func (b *bufferedNSQPublisher) Publish(message Message) error {
	b.messageCh <- message
	return nil
}

func (b *bufferedNSQPublisher) startLoop() {
	timer := time.NewTimer(b.flushInterval)
	buf := make([]Message, b.bufferSize)
	idx := 0

	go func() {
		for {
			select {
			case message, ok := <-b.messageCh:
				if ok {
					buf[idx] = message
					idx++
				}

				if !ok || idx == len(buf) {
					timer.Stop()
					b.flush(buf[:idx])
					idx = 0
					timer.Reset(b.flushInterval)
				}
			case <-timer.C:
				b.flush(buf[:idx])
				idx = 0
				timer.Reset(b.flushInterval)
			}
		}
	}()
}

func (b *bufferedNSQPublisher) flush(msgs []Message) {
	doneCh := make(chan *nsq.ProducerTransaction, 1)

	// create a producer transaction and if no error was returned, wait in
	// a goroutine until the transaction completes
	if err := b.producer.MultiPublishAsync(b.topicName, byts, doneCh); err != nil {
		log.Println(err)
		return
	}

	b.wg.Add(1)
	go func() {
		<-doneCh
		b.wg.Done()
	}()
}
