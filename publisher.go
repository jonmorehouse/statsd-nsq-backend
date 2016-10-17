package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
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

func (b *bufferedNSQPublisher) Stop() error {
	// closing this channel, results in a final flush occurring
	close(b.messageCh)
	b.wg.Wait()
	return nil
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
		localFlush := func() {
			timer.Stop()
			if err := b.flush(buf[:idx]); err != nil {
				log.Println("flush.error=" + err.Error())
			}
			idx = 0
			timer.Reset(b.flushInterval)
		}

		for {
			select {
			case message, ok := <-b.messageCh:
				if ok {
					buf[idx] = message
					idx++
				}
				if !ok || idx == len(buf) {
					localFlush()
				}
			case <-timer.C:
				localFlush()
			}
		}
	}()
}

func (b *bufferedNSQPublisher) flush(msgs []Message) error {
	// build out the messages JSON buffer
	messageBodies := make([][]byte, len(msgs))
	for idx, msg := range msgs {
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		messageBodies[idx] = jsonBytes
	}

	// publish the messages to NSQ via a producer transaction and when
	// finished, update the global wait group.
	doneCh := make(chan *nsq.ProducerTransaction, 1)
	if err := b.producer.MultiPublishAsync(b.topicName, messageBodies, doneCh); err != nil {
		return err
	}

	b.wg.Add(1)
	go func() {
		<-doneCh
		b.wg.Done()
	}()

	return nil
}
