package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/LdDl/license_plate_recognition"
	engine "github.com/LdDl/license_plate_recognition"
	"github.com/LdDl/license_plate_recognition/service/rpc"
)

var (
	confFile = flag.String("cfg", "conf.toml", "Path to TOML configuration file")
)

func main() {
	flag.Parse()

	cfgBytes, err := os.ReadFile(*confFile)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}
	var conf engine.Configuration
	err = toml.Unmarshal(cfgBytes, &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conf.YOLOPlates.Threshold <= 0.0 {
		conf.YOLOPlates.Threshold = 0.2
	}
	if conf.YOLOOCR.Threshold <= 0.0 {
		conf.YOLOOCR.Threshold = 0.3
	}
	netw, err := engine.NewYOLONetwork(conf.YOLOPlates.Cfg, conf.YOLOPlates.Weights, conf.YOLOOCR.Cfg, conf.YOLOOCR.Weights, conf.YOLOPlates.Threshold, conf.YOLOOCR.Threshold)
	if err != nil {
		fmt.Println(err)
		return
	}

	if conf.ServerConf.QueueLimit < 1 {
		conf.ServerConf.QueueLimit = 1
	}

	// Init queue worker
	q := license_plate_recognition.NewLPRQueue(netw, conf.ServerConf.QueueLimit, conf.ServerConf.SaveDetected)
	// Start worker
	q.WaitRequests()

	// Init microservice
	ms, err := rpc.NewMicroserice(q)
	if err != nil {
		fmt.Println(err)
		return
	}
	stdListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.ServerConf.Host, conf.ServerConf.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Start microservice
	if err := ms.Serve(stdListener); err != nil {
		fmt.Println(err)
		return
	}
}
