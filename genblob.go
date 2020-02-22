package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	for i := 0; i < 2; i++ {
		x0 := 100. * rand.Float64()
		y0 := 100. * rand.Float64()

		for i := 0; i < N/2; i++ {
			th := 6.28 * rand.Float64()
			r := 50. * rand.Float64()
			x := x0 + r*math.Cos(th)
			y := y0 + r*math.Sin(th)
			fmt.Printf("%f %f\n", x, y)
		}
	}
}
