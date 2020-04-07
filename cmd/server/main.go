package main

import (
	"log"
	engine "plates_recognition_grpc"
)

func main() {
	net, err := engine.NewYOLONetwork("", "", "", "")
	if err != nil {
		log.Fatalln(err)
	}
	_ = net
}
