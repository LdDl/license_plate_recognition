package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	engine "plates_recognition_grpc"
	"time"

	"google.golang.org/grpc"
)

var (
	// License plates
	platesConfig  = flag.String("platesConfig", "../data/license_plates_inference.cfg", "Path to LICENSE_PLATES network layer configuration file. Example: yolov3-plates.cfg")
	platesWeights = flag.String("platesWeights", "../data/license_plates_15000.weights", "Path to weights file. Example: yolov3-plates.weights")

	// OCR
	ocrConfig  = flag.String("ocrConfig", "../data/ocr_plates_inference.cfg", "Path to OCR network layer configuration file. Example: yolov3-ocr.cfg")
	ocrWeights = flag.String("ocrWeights", "../data/ocr_plates_7000.weights", "Path to weights file. Example: yolov3-ocr.weights")

	// gRPC port
	portConfig = flag.String("port", "50051", "Port to listen")

	// frames limit in queue
	framesLimitConfig = flag.Int("framesLimit", 200, "Max number of frames in queue")
)

func main() {
	flag.Parse()
	if *platesConfig == "" || *platesWeights == "" || *ocrConfig == "" || *ocrWeights == "" || *framesLimitConfig == 0 {
		flag.Usage()
		return
	}

	netw, err := engine.NewYOLONetwork(*platesConfig, *platesWeights, *ocrConfig, *ocrWeights)
	if err != nil {
		log.Fatalln(err)
	}

	stdListener, err := net.Listen("tcp", "0.0.0.0:"+*portConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	grpcInstance := grpc.NewServer()

	rs := &RecognitionServer{
		netW:        netw,
		framesQueue: make(chan interface{}, *framesLimitConfig),
		maxLen:      *framesLimitConfig,
	}
	rs.WaitFrames()

	engine.RegisterSTYoloServer(
		grpcInstance,
		rs,
	)

	if err := grpcInstance.Serve(stdListener); err != nil {
		log.Fatal(err)
		return
	}

}

type RecognitionServer struct {
	engine.UnimplementedSTYoloServer
	netW *engine.YOLONetwork

	framesQueue chan interface{}
	maxLen      int
}

func (rs *RecognitionServer) WaitFrames() {
	fmt.Println("YOLO networks waiting for frames now")
	go func() {
		for {
			select {
			case n := <-rs.framesQueue:
				_ = n
				_ = rs.netW // @todo распознавание
				continue
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (rs *RecognitionServer) SendToQueue(n interface{}) {
	if len(rs.framesQueue) < rs.maxLen {
		rs.framesQueue <- n
	}
}
