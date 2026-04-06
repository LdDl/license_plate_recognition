# Benchmark: go-darknet (current backend)

## Hardware
- **CPU:** Intel Core i5-10600K @ 4.10GHz
- **GPU:** NVIDIA GeForce RTX 3060
- **CUDA:** 13.1, cuDNN 9.18.1
- **OS:** Linux 6.18.7-2-cachyos

## Software
- **Go:** 1.24+
- **go-darknet:** v1.3.8
- **darknet:** AlexeyAB/darknet, GPU=1, OPENCV=0

## Models
- **Plate detection:** YOLOv4, 416x416, 6 classes
- **OCR:** YOLOv4, 416x416, 22 classes

## Test image
- `sample.jpg`: 4032x3024, 2 license plates

## Results (50 iterations, 3 warmup runs)

### Full pipeline (plate detection + OCR)

| Metric | go-darknet v1.3.8 (GPU) |
|--------|------------------------|
| Avg/frame | 155 ms |
| FPS | 6.5 |

### Plate detection only

| Metric | With Image2Float32 | Detect only (pre-converted) |
|--------|-------------------|----------------------------|
| Avg/frame | 132 ms | 33 ms |
| FPS | 7.6 | 30.3 |

### Image conversion overhead

| Method | Avg/conversion |
|--------|---------------|
| Original `Image2Float32` (v1.3.8, with `draw.Copy`) | 82 ms |
| Direct NRGBA via CGO (experimental `bench_plates_only_fixed.go`) | 89 ms |

The experimental direct-NRGBA conversion does not improve performance (0.9x). The CGO per-pixel overhead outweighs the savings from removing `draw.Copy`.

### Overhead analysis

| Component | Time |
|-----------|------|
| Image conversion (NRGBA => CHW float32) | ~82 ms |
| GPU inference | ~33 ms |
| **Total plate detection** | **~132 ms** |

Image conversion accounts for **62%** of total plate detection time.

## How to reproduce

Requires [AlexeyAB/darknet](https://github.com/AlexeyAB/darknet): `libdarknet.so` must be in the linker search path (e.g. `/usr/local/lib`) and `darknet.h` must be findable via `CGO_CFLAGS` if not in a standard include path.

```bash
cd cmd/benchmark

# If darknet.h is not in a standard include path, set CGO_CFLAGS:
# export CGO_CFLAGS="-I/path/to/darknet/include"

# Full pipeline (detect + OCR)
go run bench_darknet.go

# Plate detection only (with/without Image2Float32)
go run bench_plates_only.go

# Experimental: direct NRGBA conversion
go run bench_plates_only_fixed.go
```
