package bridge

/*
#include "od_bridge.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// RawDetection holds a single detection result from the C API.
type RawDetection struct {
	BboxX, BboxY, BboxW, BboxH int
	ClassID                    int
	Confidence                 float32
}

// Handle wraps the opaque C model pointer.
type Handle struct {
	ptr *C.struct_ModelHandle
}

// CreateModel loads an ONNX model via od-bridge (CPU).
func CreateModel(modelPath string, inputW, inputH int) (*Handle, error) {
	cPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cPath))

	ptr := C.od_model_create(cPath, C.uint32_t(inputW), C.uint32_t(inputH))
	if ptr == nil {
		return nil, fmt.Errorf("od_model_create failed for %q", modelPath)
	}
	return &Handle{ptr: ptr}, nil
}

// CreateModelCUDA loads an ONNX model via od-bridge with CUDA backend.
func CreateModelCUDA(modelPath string, inputW, inputH int) (*Handle, error) {
	cPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cPath))

	ptr := C.od_model_create_cuda(cPath, C.uint32_t(inputW), C.uint32_t(inputH))
	if ptr == nil {
		return nil, fmt.Errorf("od_model_create_cuda failed for %q", modelPath)
	}
	return &Handle{ptr: ptr}, nil
}

// Free releases the underlying Rust model.
func (h *Handle) Free() {
	if h.ptr != nil {
		C.od_model_free(h.ptr)
		h.ptr = nil
	}
}

// Detect runs inference on raw RGB bytes and returns detections.
func (h *Handle) Detect(rgb []byte, w, h2 int, confThreshold, nmsThreshold float32) ([]RawDetection, error) {
	if h.ptr == nil {
		return nil, fmt.Errorf("handle is closed")
	}

	var out C.struct_OdDetections
	rc := C.od_model_detect(
		h.ptr,
		(*C.uint8_t)(unsafe.Pointer(&rgb[0])),
		C.int32_t(w),
		C.int32_t(h2),
		C.float(confThreshold),
		C.float(nmsThreshold),
		&out,
	)

	if rc != C.Ok {
		return nil, fmt.Errorf("od_model_detect returned error code %d", int(rc))
	}

	defer C.od_detections_free(&out)

	count := int(out.len)
	if count == 0 {
		return nil, nil
	}

	cSlice := unsafe.Slice(out.data, count)
	results := make([]RawDetection, count)
	for i := 0; i < count; i++ {
		d := cSlice[i]
		results[i] = RawDetection{
			BboxX:      int(d.bbox_x),
			BboxY:      int(d.bbox_y),
			BboxW:      int(d.bbox_w),
			BboxH:      int(d.bbox_h),
			ClassID:    int(d.class_id),
			Confidence: float32(d.confidence),
		}
	}
	return results, nil
}
