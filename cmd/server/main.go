package main

import (
	"flag"
	"log"
	"net"
	engine "plates_recognition_grpc"

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
)

func main() {
	flag.Parse()
	if *platesConfig == "" || *platesWeights == "" || *ocrConfig == "" || *ocrWeights == "" {
		flag.Usage()
		return
	}

	netw, err := engine.NewYOLONetwork(*platesConfig, *platesWeights, *ocrConfig, *ocrWeights)
	if err != nil {
		log.Fatalln(err)
	}
	_ = netw

	stdListener, err := net.Listen("tcp", "0.0.0.0:"+*portConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	_ = stdListener

	grpcInstance := grpc.NewServer()

}

type RecognitionServer struct {
	engine.UnimplementedSTYoloServer
	// net chan *yolo_net.YOLONetwork
	// sync.Mutex
	// channels map[string]*yolo_net.STYolo_ConfigUpdaterServer
	// configs  map[string]*yolo_net.Config
}
