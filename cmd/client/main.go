package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/LdDl/license_plate_recognition/service/rpc/protos"

	"google.golang.org/grpc"
)

var (
	hostConfig = flag.String("host", "0.0.0.0", "server's hostname")
	portConfig = flag.String("port", "50051", "server's port")
	fileConfig = flag.String("file", "sample.jpg", "filename")

	xConfig = flag.Int("x", 0, "x (left top of crop rectangle)")
	yConfig = flag.Int("y", 0, "y (left top of crop rectangle)")

	widthConfig  = flag.Int("width", 4032, "width of crop rectangle")
	heightConfig = flag.Int("height", 3024, "height of crop rectangle")
)

func main() {
	flag.Parse()

	if *hostConfig == "" || *portConfig == "" || *fileConfig == "" {
		flag.Usage()
		return
	}

	// Read image from file
	ifile, err := os.Open(*fileConfig)
	if err != nil {
		log.Println(err)
		return
	}
	imgIn, _, err := image.Decode(ifile)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imgIn, nil)
	if err != nil {
		log.Println(err)
		return
	}
	sendS3 := buf.Bytes()

	// Connect to gRPC
	url := fmt.Sprintf("%s:%s", *hostConfig, *portConfig)
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Init gRPC client
	client := protos.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Send message to gRPC server
	resp, err := client.ProcessImage(
		ctx,
		&protos.LPRRequest{
			Image: sendS3,
			Bbox: &protos.BBox{
				XLeft:  int32(*xConfig),
				YTop:   int32(*yConfig),
				Width:  int32(*widthConfig),
				Height: int32(*heightConfig),
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Elapsed seconds:", resp.Elapsed)
	fmt.Println("Detections num:", len(resp.LicensePlates))
	for i, detection := range resp.LicensePlates {
		fmt.Printf("Detection #%d:\n", i)
		fmt.Println("\tText:", detection.Text)
		fmt.Println("\tPlate bbox:", detection.Bbox)
		fmt.Println("\tOCR bboxes:", detection.OcrBboxes)
	}
}
