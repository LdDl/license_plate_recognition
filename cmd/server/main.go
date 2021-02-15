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
	"time"

	engine "github.com/LdDl/license_plate_recognition"
	grpc_server "github.com/LdDl/odam"
	"github.com/disintegration/imaging"
	"google.golang.org/grpc"
)

var (
	// License plates
	platesConfig  = flag.String("platesConfig", "../data/license_plates_inference.cfg", "Path to LICENSE_PLATES network layer configuration file. Example: yolov3-plates.cfg")
	platesWeights = flag.String("platesWeights", "../data/license_plates_15000.weights", "Path to weights file. Example: yolov3-plates.weights")

	// OCR
	ocrConfig  = flag.String("ocrConfig", "../data/ocr_plates_inference.cfg", "Path to OCR network layer configuration file. Example: yolov3-ocr.cfg")
	ocrWeights = flag.String("ocrWeights", "../data/ocr_plates_140000.weights", "Path to weights file. Example: yolov3-ocr.weights")

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
		netW:          netw,
		framesQueue:   make(chan *vehicleInfo, *framesLimitConfig),
		maxLen:        *framesLimitConfig,
		resp:          make(chan *ServerResponse, *framesLimitConfig),
		AfterFunction: doSomeStuff,
	}
	// Init neural network's queue
	rs.WaitFrames()

	// Init gRPC server
	grpcInstance := grpc.NewServer()

	// Register servers
	grpc_server.RegisterServiceYOLOServer(
		grpcInstance,
		rs,
	)

	// Start
	if err := grpcInstance.Serve(stdListener); err != nil {
		log.Fatal(err)
		return
	}

}

// RecognitionServer Wrapper around engine.ServiceYOLOServer
type RecognitionServer struct {
	grpc_server.ServiceYOLOServer
	netW        *engine.YOLONetwork
	framesQueue chan *vehicleInfo
	maxLen      int
	resp        chan *ServerResponse

	AfterFunction func(data *PlateInfo, fileContents []byte) error
}

// ServerResponse Response from server
type ServerResponse struct {
	Resp  *engine.YOLOResponse
	Error error
}

// WaitFrames Endless loop for waiting frames
func (rs *RecognitionServer) WaitFrames() {
	fmt.Println("YOLO networks waiting for frames now")
	go func() {
		for {
			select {
			case n := <-rs.framesQueue:
				// fmt.Println("img of size", n.Bounds().Dx(), n.Bounds().Dy())
				resp, err := rs.netW.ReadLicensePlates(n.img, true)

				if *saveDetectedConfig != 0 {
					for i := range resp.Plates {
						err := ensureDir("./detected")
						if err != nil {
							fmt.Println("Can't check or create directory './detected':", err)
							continue
						}
						fname := fmt.Sprintf("./detected/%s_%s_%.0f.jpeg", resp.Plates[i].Text, time.Now().Format("2006-01-02T15-04-05"), resp.Plates[i].Probability)
						f, err := os.Create(fname)
						if err != nil {
							fmt.Println("Can't create file:", err)
							continue
						}
						defer f.Close()

						err = jpeg.Encode(f, resp.Plates[i].CroppedNumber, nil)
						if err != nil {
							fmt.Println("Can't encode JPEG:", err)
							continue
						}

						if resp.Plates[i].Text != "" {
							dplat := PlateInfo{
								CameraID: n.imageInfo.CamId,
								Text:     resp.Plates[i].Text,
								Time:     time.Now().UTC().Format("2006-01-02T15:04:05"),
							}
							copyBuff := new(bytes.Buffer)
							err = jpeg.Encode(copyBuff, resp.Plates[i].CroppedNumber, nil)
							fileContents := copyBuff.Bytes()
							err := rs.AfterFunction(&dplat, fileContents)
							if err != nil {
								fmt.Println("Can't exectude AfterFunction:", err)
								continue
							}
						}

					}
				}

				rs.resp <- &ServerResponse{resp, err}
				continue
			}
		}
	}()
}

// SendToQueue Add element to queue
func (rs *RecognitionServer) SendToQueue(n *vehicleInfo) {
	if len(rs.framesQueue) < rs.maxLen {
		rs.framesQueue <- n
	}
}

type vehicleInfo struct {
	imageInfo *grpc_server.ObjectInformation
	img       *image.NRGBA
}

// SendDetection Imeplented function or accepting image
func (rs *RecognitionServer) SendDetection(ctx context.Context, in *grpc_server.ObjectInformation) (*grpc_server.Response, error) {

	imgBytes := in.GetImage()
	imgReader := bytes.NewReader(imgBytes)

	stdImage, _, err := image.Decode(imgReader)
	if err != nil {
		return &grpc_server.Response{Error: "Image decoding failed"}, err
	}

	height := stdImage.Bounds().Dy()
	width := stdImage.Bounds().Dx()

	det := in.GetDetection()
	xl := int(det.GetXLeft())
	yt := int(det.GetYTop())
	dh := int(det.GetHeight())
	dw := int(det.GetWidth())
	if dw <= 0 || dh <= 0 || xl >= width || yt >= height {
		return &grpc_server.Response{Error: "Incorrect bounding box of a car"}, nil
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

	inf := vehicleInfo{
		imageInfo: in,
		img:       vehicleImg,
	}
	rs.SendToQueue(&inf)

	response := <-rs.resp
	log.Println(response.Resp.String())
	if response.Error != nil {
		return &grpc_server.Response{Message: "error", Warning: response.Error.Error()}, nil
	}
	return &grpc_server.Response{Message: "ok", Warning: ""}, nil
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0777)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}

// PlateInfo Information about license plate
type PlateInfo struct {
	CameraID string `json:"camera_id" example:"f2abe45e-aad8-40a2-a3b7-0c610c0f3dda"`
	Text     string `json:"text" example:"a777aa77"`
	Time     string `json:"tm" example:"2020-04-30T00:00:00"`
}

func doSomeStuff(data *PlateInfo, fileContents []byte) error {
	/*
		If you want, you can implement this function by yourself (and you can wrap this function also)
		Default behaviour: do nothing.
	*/
	return nil
}
