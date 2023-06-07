package rpc

import (
	"github.com/LdDl/license_plate_recognition"
	"github.com/LdDl/license_plate_recognition/service/rpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Microservice struct {
	protos.ServiceServer
	engine *license_plate_recognition.LPRQueue
}

func NewMicroserice(engine *license_plate_recognition.LPRQueue) (*grpc.Server, error) {
	grpcInstance := grpc.NewServer()
	server := Microservice{
		engine: engine,
	}
	protos.RegisterServiceServer(
		grpcInstance,
		&server,
	)
	reflection.Register(grpcInstance)
	return grpcInstance, nil
}

// microservice is just an alias to internal implementation
func (microservice *Microservice) SendToQueue(req *license_plate_recognition.QueueRequest) (*license_plate_recognition.QueueResponse, error) {
	return microservice.engine.SendToQueue(req)
}
