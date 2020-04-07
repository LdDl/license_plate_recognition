package main

var (
	hostConfig = flag.String("host", "0.0.0.0", "server's hostname")
	portConfig = flag.String("port", "50051", "server's port")
	fileConfig = flag.String("file", "sample.jpg", "filename")

	xConfig = flag.String("x", "0", "x (left top of crop rectangle)")
	yConfig = flag.String("y", "0", "y (left top of crop rectangle)")

	widthConfig  = flag.String("width", "4032", "width of crop rectangle")
	heightConfig = flag.String("height", "3024", "height of crop rectangle")
)


func main() {

}