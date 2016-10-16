package main

type Message interface{}

type Publisher interface {
	Publish(Message) error
}

func NewNSQPublisher() Publisher {
	return nil
}

// nsqPublisher is a publisher which emits messages to NSQ for each metric it
type nsqPublisher struct {
	//
}

func NewSampledPublisher(publisher Publisher, sampleRate float64) Publisher {
	return nil
}

type sampledPublisher struct {
	//
}

func NewBufferedPublisher(publisher Publisher, bufferSize int) Publisher {
	return &bufferedPublisher{
		buffer: make([]Message, 0, bufferSize),
	}
}

type bufferedPublisher struct {
	publisher Publisher
	buffer    []Message
}

func (b *bufferedPublisher) Publish(Message) error {
	return nil
}
