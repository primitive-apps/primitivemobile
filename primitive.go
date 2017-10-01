package primitivemobile

import (
	"log"
	"math/rand"
	"runtime"
	"time"
	"os"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

/*var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray 
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V, VV      bool
)*/


/*
func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}
*/


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessImage(Input string, InputSize int, OutputSize int, Count int, Mode int, Background string, Alpha int, Repeat int, Output string) {
/*	Input := fileName
	InputSize := 256
	OutputSize := 1024
	Count := 1
	Mode := 1
	Background := ""
	Alpha := 128
	Repeat := 0*/



	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// determine worker count
	Workers := runtime.NumCPU()

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}

	// run algorithm
	model := primitive.NewModel(input, bg, OutputSize, Workers)
	frame := 0

	outputFile, err := os.OpenFile(Output, os.O_WRONLY, 0666)
	if err != nil{
		log.Fatal(err)
	}

	for i := 0; i < Count; i++ {
		frame++

		// find optimal shape and add it to the model
		//t := time.Now()

		model.Step(primitive.ShapeType(Mode), Alpha, Repeat)
		outputFile.Seek(0,0)
		outputFile.WriteString(model.SVG())
		outputFile.Sync()
		
/* 		nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
		elapsed := time.Since(start).Seconds()
		primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps) */

		// write output image(s)
	}
	//return model.SVG();
	outputFile.Close()
}
