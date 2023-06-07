package license_plate_recognition

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"
)

var (
	ErrMaxTimeout = fmt.Errorf("max timeout is 1 minute")
)

// LPRQueue is a queue to handle asynchronous requests to neural network
type LPRQueue struct {
	network       *YOLONetwork
	requestsQueue chan *QueueRequest
	queueLimit    int
	saveDetected  bool
}

func NewLPRQueue(network *YOLONetwork, queueLimit int, saveDetected bool) *LPRQueue {
	return &LPRQueue{
		network:       network,
		requestsQueue: make(chan *QueueRequest, queueLimit),
		queueLimit:    queueLimit,
		saveDetected:  saveDetected,
	}
}

// QueueRequest is wrapping around image and response channel
type QueueRequest struct {
	ImageData    *image.NRGBA
	responseChan chan *QueueResponse
	ctx          context.Context
}

func NewQueueRequest(ctx context.Context, image *image.NRGBA) QueueRequest {
	if ctx == nil {
		return QueueRequest{
			ImageData: image,
			ctx:       context.Background(),
		}
	}
	return QueueRequest{
		ImageData: image,
		ctx:       ctx,
	}
}

// QueueResponse is just response from YOLO
type QueueResponse struct {
	Resp  *YOLOResponse
	Error error
}

// SendToQueue is wrapper around request and its context
func (q *LPRQueue) SendToQueue(req *QueueRequest) (*QueueResponse, error) {
	if len(q.requestsQueue) >= q.queueLimit {
		return nil, fmt.Errorf("queue is full, unable to send request")
	}
	req.responseChan = make(chan *QueueResponse)
	q.requestsQueue <- req
	select {
	case resp := <-req.responseChan:
		return resp, nil
	case <-req.ctx.Done():
		return nil, req.ctx.Err()
	case <-time.After(60 * time.Second):
		return nil, ErrMaxTimeout
	}
}

// WaitRequests is endless loop for waiting frames
func (q *LPRQueue) WaitRequests() {
	fmt.Println("YOLO networks waiting for requests...")
	go func() {
		for req := range q.requestsQueue {
			resp, err := q.network.ReadLicensePlates(req.ImageData, true)
			if err != nil {
				req.responseChan <- &QueueResponse{resp, err}
				continue
			}
			if q.saveDetected {
				for _, plate := range resp.Plates {
					err := ensureDir("./detected")
					if err != nil {
						fmt.Println("Can't check or create directory './detected':", err)
						continue
					}
					tm := time.Now().Format("2006-01-02T15-04-05")
					imageFileName := fmt.Sprintf("./detected/%s_%s_%.0f.jpeg", plate.Text, tm, plate.Probability)
					file, err := os.Create(imageFileName)
					if err != nil {
						req.responseChan <- &QueueResponse{resp, err}
						continue
					}
					err = jpeg.Encode(file, plate.CroppedNumber, nil)
					if err != nil {
						file.Close()
						req.responseChan <- &QueueResponse{resp, err}
						continue
					}
					file.Close() // Explicit close is needed in infitite loop

					annotationFileName := fmt.Sprintf("./detected/%s_%s_%.0f.txt", plate.Text, tm, plate.Probability)
					preparedAnnotations := []string{}
					width := float64(plate.CroppedNumber.Rect.Dx())
					height := float64(plate.CroppedNumber.Rect.Dy())
					for charIdx, char := range plate.OCRRects {
						yoloAnn := pascalVOC2YOLO(
							float64(char.Min.X), float64(char.Min.Y),
							float64(char.Max.X), float64(char.Max.Y),
							width, height,
						)
						classID := plate.OCRClassesIDs[charIdx]
						preparedAnnotation := []string{
							fmt.Sprintf("%d", classID),
							fmt.Sprintf("%.24f", yoloAnn[0]),
							fmt.Sprintf("%.24f", yoloAnn[1]),
							fmt.Sprintf("%.24f", yoloAnn[2]),
							fmt.Sprintf("%.24f", yoloAnn[3]),
						}
						preparedAnnotations = append(preparedAnnotations, strings.Join(preparedAnnotation, " "))
					}
					annFile, err := os.Create(annotationFileName)
					if err != nil {
						req.responseChan <- &QueueResponse{resp, err}
						continue
					}
					_, err = annFile.WriteString(strings.Join(preparedAnnotations, "\n"))
					if err != nil {
						annFile.Close()
						req.responseChan <- &QueueResponse{resp, err}
						continue
					}
					annFile.Close()
				}
			}
			req.responseChan <- &QueueResponse{resp, err}
			continue
		}
	}()
}
