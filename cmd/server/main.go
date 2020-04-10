package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net"
	"os"
	engine "plates_recognition_grpc"
	"time"

	"github.com/disintegration/imaging"
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

	// Store detected object to filesytem?
	saveDetectedConfig = flag.Uint("saveDetected", 0, "Do you want to save detected objects into JPEG files?")
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

	// Init servers
	rs := &RecognitionServer{
		netW:        netw,
		framesQueue: make(chan *image.NRGBA, *framesLimitConfig),
		maxLen:      *framesLimitConfig,
		resp:        make(chan *ServerResponse, *framesLimitConfig),
	}
	// Init neural network's queue
	rs.WaitFrames()

	// Init gRPC server
	grpcInstance := grpc.NewServer()

	// Register servers
	engine.RegisterSTYoloServer(
		grpcInstance,
		rs,
	)

	// Start
	if err := grpcInstance.Serve(stdListener); err != nil {
		log.Fatal(err)
		return
	}

}

type RecognitionServer struct {
	engine.UnimplementedSTYoloServer
	netW *engine.YOLONetwork

	framesQueue chan *image.NRGBA
	maxLen      int

	resp chan *ServerResponse
}

type ServerResponse struct {
	Resp  *engine.YOLOResponse
	Error error
}

func (rs *RecognitionServer) WaitFrames() {
	fmt.Println("YOLO networks waiting for frames now")
	go func() {
		for {
			select {
			case n := <-rs.framesQueue:
				// fmt.Println("img of size", n.Bounds().Dx(), n.Bounds().Dy())
				resp, err := rs.netW.ReadLicensePlates(n, true)

				if *saveDetectedConfig != 0 {
					for i := range resp.Plates {
						fname := fmt.Sprintf("./detected/%s_%s_%.0f.jpeg", resp.Plates[i].Text, time.Now().Format("2006-01-02T15-04-05"), resp.Plates[i].Probability)
						f, err := os.Create(fname)
						if err != nil {
							fmt.Println(err)
							// rs.resp <- &ServerResponse{nil, err}
						}
						defer f.Close()

						err = jpeg.Encode(f, resp.Plates[i].CroppedNumber, nil)
						if err != nil {
							fmt.Println(err)
							// rs.resp <- &ServerResponse{nil, err}
						}

					}
				}

				rs.resp <- &ServerResponse{resp, err}
				continue
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (rs *RecognitionServer) SendToQueue(n *image.NRGBA) {
	if len(rs.framesQueue) < rs.maxLen {
		rs.framesQueue <- n
	}
}

func (rs *RecognitionServer) SendDetection(ctx context.Context, in *engine.CamInfo) (*engine.Response, error) {
	imgBytes := in.GetImage()
	imgReader := bytes.NewReader(imgBytes)

	stdImage, _, err := image.Decode(imgReader)
	if err != nil {
		return &engine.Response{Error: "Image decoding failed"}, err
	}

	height := stdImage.Bounds().Dy()
	width := stdImage.Bounds().Dx()

	det := in.GetDetection()
	xl := int(det.GetXLeft())
	yt := int(det.GetYTop())
	dh := int(det.GetHeight())
	dw := int(det.GetWidth())
	if dw <= 0 || dh <= 0 || xl >= width || yt >= height {
		return &engine.Response{Error: "Incorrect bounding box of a car"}, nil
	}

	bbw := xl + dw
	bbh := yt + dh
	if xl < 0 || yt < 0 || xl+dw > width || yt+dh > height {
		// Bounding box is bigger than image
	}
	if xl < 0 {
		xl = 0
	}
	if yt < 0 {
		yt = 0
	}
	if bbw > width {
		bbw = width
	}
	if bbh > height {
		bbh = height
	}

	vehicleBBox := image.Rect(xl, yt, bbw, bbh)
	vehicleImg := imaging.Crop(stdImage, vehicleBBox)

	rs.SendToQueue(vehicleImg)

	response := <-rs.resp
	log.Println(response.Resp.String())
	if response.Error != nil {
		return &engine.Response{Message: "error", Warning: response.Error.Error()}, nil
	}
	return &engine.Response{Message: "ok", Warning: ""}, nil
}
