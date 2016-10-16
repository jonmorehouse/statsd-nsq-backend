package main

import (
	"log"
	"strconv"
	"time"
)

const OptimalPayloadSize = 1432

/*
MaxUDPPayloadSize defines the maximum payload size for a UDP datagram.
Its value comes from the calculation: 65535 bytes Max UDP datagram size -
8byte UDP header - 60byte max IP headers
any number greater than that will see frames being cut out.
*/
const MaxUDPPayloadSize = 65467

func main() {
	ch := make(chan []byte)
	go func(ch chan []byte) {
		for {
			byts, ok := <-ch
			if !ok {
				return
			}
			log.Println(string(byts))
		}
	}(ch)

	server := NewServer(time.Second*10, ch)
	if err := server.ListenAnywhere(); err != nil {
		log.Fatal("unable to start server. err=" + err.Error())
	}

	log.Println("listening on port: " + strconv.Itoa(server.Port()))
	time.Sleep(time.Minute * 10)
	if err := server.Stop(); err != nil {
		log.Println("error received: " + err.Error())
	}

	close(ch)
}
