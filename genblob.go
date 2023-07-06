package main

/*
 * genblob - write text coordinates of circular "blobs"
 * or clusters of points on stdout.
 *
 * Usage: genblob $clusters $total_points
 */

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
	max, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	N, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	// max is number of "blobs" or clusters
	for i := 0; i < max; i++ {
		x0 := 100. * rand.Float64()
		y0 := 100. * rand.Float64()
		fmt.Fprintf(os.Stderr, "# Blob %d X0,Y0 %f,%f\n", i, x0, y0)

		// N/max count of points per blob.
		// This will produce point-density higher the closter
		// to <x0,y0> you get, since a given "racetrack"
		// has more area as radius increases, but the uniform
		// distribution of rand.Float64 will put the same number
		// of points on any given radius.
		for i := 0; i < N/max; i++ {
			th := 6.28 * rand.Float64()
			r := 50. * rand.Float64()
			x := x0 + r*math.Cos(th)
			y := y0 + r*math.Sin(th)
			fmt.Printf("%f %f\n", x, y)
		}
	}
}
