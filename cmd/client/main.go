package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	engine "github.com/LdDl/odam"

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
	client := engine.NewServiceYOLOClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Send message to gRPC server
	r, err := client.SendDetection(
		ctx,
		&engine.ObjectInformation{
			CamId:     "my_new_uuid",
			Timestamp: time.Now().Unix(),
			Image:     sendS3,
			Class: &engine.ClassInfo{
				ClassId:   100,
				ClassName: "find_ocr",
			},
			Detection: &engine.Detection{
				XLeft:  int32(*xConfig),
				YTop:   int32(*yConfig),
				Width:  int32(*widthConfig),
				Height: int32(*heightConfig),
			},
			// Skip virtual line part (not needed)
			VirtualLine: &engine.VirtualLineInfo{},
			// Skip tracking info part (not needed)
			TrackInformation: &engine.TrackInfo{},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	if len(r.GetError()) != 0 {
		log.Fatalln(r.GetError())
	}

	if len(r.GetWarning()) != 0 {
		log.Println("Warn:", r.GetWarning())
	}

	log.Println("Answer:", r.GetMessage())
}
