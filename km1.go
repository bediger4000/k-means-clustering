package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// Point an x,y cartesian point
type Point struct {
	x float64
	y float64
}

func main() {
	k, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	filename := os.Args[1]
	points := readPoints(filename)

	ps := PointSlice(points)
	sort.Sort(ps)

	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	kmeanscluster(k, points)
}

func kmeanscluster(k int, points []Point) {

	var centroids []Point

	for len(centroids) < k {
		candidate := points[rand.Intn(len(points))]
		foundit := false
		for _, centroid := range centroids {
			if candidate.x == centroid.x && candidate.y == centroid.y {
				foundit = true
				break
			}
		}
		if !foundit {
			centroids = append(centroids, candidate)
		}
	}

	looping := true

	var finalclusters [][]Point

	for looping {
		// fmt.Printf("Centroids: %+v\n", centroids)

		clusters := make([][]Point, k)

		for _, point := range points {
			distances := make([]float64, k)
			for i := 0; i < k; i++ {
				centroid := centroids[i]
				dx := centroid.x - point.x
				dy := centroid.y - point.y
				distances[i] = dx*dx + dy*dy
			}
			min := distances[0]
			cent := 0
			for i, dist := range distances {
				if dist < min {
					min = dist
					cent = i
				}
			}
			clusters[cent] = append(clusters[cent], point)
		}

		newcentroids := calcCentroids(clusters)
		looping = compareCentroids(centroids, newcentroids)
		centroids = newcentroids
		finalclusters = clusters
	}

	for i, centroid := range centroids {
		fmt.Printf("%f %f c%d\n", centroid.x, centroid.y, i)
	}

	for i, cluster := range finalclusters {
		for _, point := range cluster {
			fmt.Printf("%f %f %d\n", point.x, point.y, i)
		}
	}
}

func compareCentroids(centroids []Point, newcentroids []Point) bool {
	if len(centroids) != len(newcentroids) {
		log.Fatalf("%d old centroids, %d newcentroids\n", len(centroids), len(newcentroids))
	}

	for i := 0; i < len(centroids); i++ {
		dx := centroids[i].x - newcentroids[i].x
		dy := centroids[i].y - newcentroids[i].y
		dist := dx*dx + dy*dy
		if dist > 0.01 {
			return true // keep looping
		}
	}

	return false // stop looping
}

type PointSlice []Point

func (ps PointSlice) Len() int { return len(ps) }
func (ps PointSlice) Less(i, j int) bool {
	if ps[i].x < ps[j].x {
		return true
	}
	if ps[i].x == ps[j].x {
		return ps[i].y < ps[j].y
	}
	return false
}
func (ps PointSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

func calcCentroids(clusters [][]Point) []Point {
	centroids := make([]Point, len(clusters))

	for cent, cluster := range clusters {
		var sumx, sumy float64
		for _, point := range cluster {
			sumx += point.x
			sumy += point.y
		}
		centroids[cent] = Point{x: sumx / float64(len(cluster)), y: sumy / float64(len(cluster))}
	}

	return centroids
}

func readPoints(filename string) []Point {
	fin, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	var points []Point

	for {
		var p Point
		n, err := fmt.Fscanf(fin, "%f %f\n", &p.x, &p.y)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Print(err)
			continue
		}
		if n != 2 {
			log.Printf("Parsed %d items, wanted 2\n", n)
			continue
		}

		points = append(points, p)
	}

	return points
}
