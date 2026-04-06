//go:build nopkgconfig

package bridge

/*
#cgo LDFLAGS: -lod_bridge -lm -ldl -lpthread
#cgo CFLAGS: -I/usr/local/include/od-bridge
*/
import "C"
