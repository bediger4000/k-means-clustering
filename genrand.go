package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	maxPop := flag.Int("p", 0, "maxium population")
	flag.Parse()

	args := flag.Args()

	N, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	population := *maxPop > 0

	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	if population {
		for i := 0; i < N; i++ {
			fmt.Printf("%d %f %f\n", rand.Intn(*maxPop), 1500.*rand.Float64(), 1500.0*rand.Float64())
		}
	} else {
		for i := 0; i < N; i++ {
			fmt.Printf("%f %f\n", 1500.*rand.Float64(), 1500.0*rand.Float64())
		}
	}
}
