# Benchmark: go-darknet vs od-bridge (ORT)

Baseline darknet results measured on [v1.3.5](https://github.com/LdDl/license_plate_recognition/releases/tag/v1.3.5) (last darknet release).

## Hardware
- **CPU:** Intel Core i5-10600K @ 4.10GHz
- **GPU:** NVIDIA GeForce RTX 3060
- **CUDA:** 13.1, cuDNN 9.18.1
- **OS:** Linux 6.18.7-2-cachyos

## Software
- **Go:** 1.24+
- **go-darknet:** v1.3.8 (last darknet release: [v1.3.5](https://github.com/LdDl/license_plate_recognition/releases/tag/v1.3.5))
- **od-bridge:** v0.1.0 (od_opencv + ORT 2.0-rc12, CUDA EP)
- **darknet:** AlexeyAB/darknet, GPU=1, OPENCV=0

## Models
- **Plate detection:** YOLOv4, 416x416, 6 classes
- **OCR:** YOLOv4, 416x416, 22 classes

## Test image
- `sample.jpg`: 4032x3024, 2 license plates

## Results (50 iterations, 3 warmup runs)

### Summary

| Test | go-darknet v1.3.5 (GPU) | od-bridge ORT (CUDA) | Speedup |
|------|------------------------|---------------------|---------|
| Plate detection | 132 ms / 7.6 FPS | **24 ms / 41.1 FPS** | **5.5x** |
| Full pipeline (detect + OCR) | 155 ms / 6.5 FPS | **77 ms / 13.0 FPS** | **2.0x** |

### od-bridge ORT (CUDA) details

| Test | Avg/frame | FPS |
|------|-----------|-----|
| Plate detection | 24 ms | 41.1 |
| OCR (single crop 639x303) | 20 ms | 50.3 |
| Full pipeline (detect + OCR all) | 77 ms | 13.0 |

### go-darknet v1.3.5 details

| Test | Avg/frame | FPS |
|------|-----------|-----|
| Plate detection (with Image2Float32) | 132 ms | 7.6 |
| Plate detection (Detect only, pre-converted) | 33 ms | 30.3 |
| Full pipeline (detect + OCR) | 155 ms | 6.5 |

### Overhead analysis (go-darknet)

| Component | Time |
|-----------|------|
| Image conversion (NRGBA => CHW float32) | ~82 ms |
| GPU inference | ~33 ms |
| **Total plate detection** | **~132 ms** |

Image conversion accounts for 62% of total plate detection time in go-darknet.
ORT accepts raw uint8 HWC input and fuses conversion into the GPU pipeline, eliminating this bottleneck.

## How to reproduce

### od-bridge (current)

Requires [od-bridge](https://github.com/LdDl/od-bridge) installed with CUDA support (see [installation guide](https://github.com/LdDl/od-bridge#installation)).

```bash
cd cmd/benchmark
go run bench_od_bridge.go
```

### go-darknet (baseline, v1.3.5)

Requires [AlexeyAB/darknet](https://github.com/AlexeyAB/darknet): `libdarknet.so` must be in the linker search path and `darknet.h` must be findable via `CGO_CFLAGS` if not in a standard include path.

```bash
git checkout v1.3.5
cd cmd/benchmark

# If darknet.h is not in a standard include path, set CGO_CFLAGS:
# export CGO_CFLAGS="-I/path/to/darknet/include"

go run bench_darknet.go
go run bench_plates_only.go
```
