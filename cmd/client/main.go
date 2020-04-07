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
	engine "plates_recognition_grpc"
	"time"

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

	url := fmt.Sprintf("%s:%s", *hostConfig, *portConfig)

	// Set up a connection to the server.
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	c := engine.NewSTYoloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	channel, err := c.ConfigUpdater(context.Background())
	channel.Send(&engine.Response{Message: "Channel opened!"})
	channel.Send(&engine.Response{Message: "Channel opened!"})
	cfg := &engine.Config{}
	cfg, err = channel.Recv()
	cfg.DetectionLines = []*engine.DetectionLine{&engine.DetectionLine{Id: 1, Begin: &engine.Point{X: 1, Y: 1}, End: &engine.Point{X: 416, Y: 416}}}
	resp, _ := c.SetConfig(ctx, cfg)
	fmt.Println(resp)
	defer cancel()
	r, err := c.SendDetection(
		ctx,
		&engine.CamInfo{
			CamId:     cfg.GetUid(),
			Timestamp: time.Now().Unix(),
			Image:     sendS3,
			Detection: &engine.Detection{
				XLeft:  int32(*xConfig),
				YTop:   int32(*yConfig),
				Width:  int32(*widthConfig),
				Height: int32(*heightConfig),
				LineId: 1,
			},
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	if len(r.GetError()) != 0 {
		log.Println(r.GetError())
		return
	}
	if len(r.GetWarning()) != 0 {
		log.Println("Warn:", r.GetWarning())
	}
	c.SetConfig(ctx, cfg)
	log.Println("Answer:", r.GetMessage())
}
