package main

import (
	"log"
	"strconv"
	"time"
)

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
