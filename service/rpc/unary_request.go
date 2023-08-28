package rpc

import (
	"bytes"
	"context"
	"fmt"
	"image"

	"github.com/LdDl/license_plate_recognition"
	"github.com/LdDl/license_plate_recognition/service/rpc/protos"
	"github.com/disintegration/imaging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrBadBoundingBBox = fmt.Errorf("bad bounding box")
	ErrBadDetections   = fmt.Errorf("bad detections")
)

// ProcessImage is to match gRPC server interface. Provides business logic for processing image
func (ts *Microservice) ProcessImage(ctx context.Context, in *protos.LPRRequest) (*protos.LPRResponse, error) {

	imgBytes := in.GetImage()
	imgReader := bytes.NewReader(imgBytes)

	stdImage, _, err := image.Decode(imgReader)
	if err != nil {
		return &protos.LPRResponse{Error: "Image decoding failed"}, status.Error(codes.Internal, err.Error())
	}

	height := stdImage.Bounds().Dy()
	width := stdImage.Bounds().Dx()

	dw := width
	dh := height
	xl := 0
	yt := 0

	if in.Bbox != nil {
		xl = int(in.Bbox.XLeft)
		yt = int(in.Bbox.YTop)
		dw = int(in.Bbox.Width)
		dh = int(in.Bbox.Height)
	}
	if dw <= 0 || dh <= 0 || xl >= width || yt >= height {
		return &protos.LPRResponse{Error: "Provided bounding box is incorrect"}, status.Error(codes.InvalidArgument, ErrBadBoundingBBox.Error())
	}

	bbw := xl + dw
	bbh := yt + dh
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

	customBBox := image.Rect(xl, yt, bbw, bbh)
	bboxCrop := imaging.Crop(stdImage, customBBox)

	req := license_plate_recognition.NewQueueRequest(ctx, bboxCrop)
	resp, err := ts.SendToQueue(&req)
	if err != nil {
		return &protos.LPRResponse{Error: "Can't complete job"}, status.Error(codes.Internal, err.Error())
	}
	if resp.Error != nil {
		return &protos.LPRResponse{Error: "Inference error"}, status.Error(codes.Internal, resp.Error.Error())
	}
	if resp.Resp == nil {
		return &protos.LPRResponse{Error: "Empty inference result"}, status.Error(codes.Internal, ErrBadDetections.Error())
	}
	ans := protos.LPRResponse{
		LicensePlates: make([]*protos.LPRInfo, len(resp.Resp.Plates)),
		Elapsed:       float32(resp.Resp.Elapsed.Seconds()),
		Message:       "ok",
		Warning:       "",
		Error:         "",
	}
	for i, plate := range resp.Resp.Plates {
		ans.LicensePlates[i] = &protos.LPRInfo{
			Bbox: &protos.BBox{
				XLeft:  int32(plate.Rect.Min.X),
				YTop:   int32(plate.Rect.Min.X),
				Height: int32(plate.Rect.Dy()),
				Width:  int32(plate.Rect.Dx()),
			},
			OcrBboxes: make([]*protos.BBox, len(plate.OCRRects)),
			Text:      plate.Text,
		}
		for j, char := range plate.OCRRects {
			ans.LicensePlates[i].OcrBboxes[j] = &protos.BBox{
				XLeft:  int32(char.Min.X),
				YTop:   int32(char.Min.X),
				Height: int32(char.Dy()),
				Width:  int32(char.Dx()),
			}
		}
	}
	return &ans, nil
}
