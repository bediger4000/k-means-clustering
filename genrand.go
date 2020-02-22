package main

import (
	"fmt"
	"log"
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

	for i := 0; i < N; i++ {
		fmt.Printf("%f %f\n", 100.*rand.Float64(), 100.0*rand.Float64())
	}
}
