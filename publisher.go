package main

type Publisher interface {
	Publish(Message) error
}

func NewNSQPublisher() {
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

}

type bufferedPublisher struct {
	publisher
}
