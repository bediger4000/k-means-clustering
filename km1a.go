package main

/*
   K-means clustering with k-means++ initial centroid choice.
*/

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

type dist struct {
	D2         float64
	pointIndex int
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

	centroids := kMeansPPCentroids(k, points)

	looping := true

	var finalclusters [][]Point

	for looping {

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
		var dummy float64
		n, err := fmt.Fscanf(fin, "%f %f %f\n", &dummy, &p.x, &p.y)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Print(err)
			continue
		}
		if n != 3 {
			log.Printf("Parsed %d items, wanted 3\n", n)
			continue
		}

		points = append(points, p)
	}

	return points
}

/* k-means++ method of finding initial guesses at centroids.

1. Choose one center uniformly at random among the data points.
2. For each data point x, compute D(x), the distance between x and the nearest
   center that has already been chosen.
3. Choose one new data point at random as a new center, using a weighted
   probability distribution where a point x is chosen with probability
   proportional to D(x)^2.
4. Repeat Steps 2 and 3 until k centers have been chosen.
5. Now that the initial centers have been chosen, proceed using standard k-means clustering.
*/
func kMeansPPCentroids(k int, points []Point) (centroids []Point) {

	centroids = append(centroids, points[rand.Intn(len(points))])

	D := make([]dist, len(points))

	for i := 0; i < k-1; i++ {
		fillDistances(D, points, centroids)
		newCenterIndex := weightedChoice(D)
		centroids = append(centroids, points[newCenterIndex])
	}

	return
}

type interval struct {
	maxval float64
	index  int
}

/*
   Choose one new data point at random as a new center, using a weighted
   probability distribution where a point x is chosen with probability
   proportional to D(x)^2.
*/
func weightedChoice(D []dist) int {
	sumValues := 0.0
	intervals := make([]interval, len(D))

	for i := range D {
		sumValues += D[i].D2
		intervals = append(intervals, interval{maxval: sumValues, index: D[i].pointIndex})
	}

	inInterval := rand.Float64() * sumValues

	for i := range intervals {
		if inInterval < intervals[i].maxval {
			return intervals[i].index
		}
	}

	return 0
}

/*
    For each data point x, compute D(x), the distance between x and the nearest
    center that has already been chosen.

	Actually going to calculate D(x)^2, because that's what's used later in the
	algorithm.
*/
func fillDistances(D []dist, points []Point, centroids []Point) {

	for idx, point := range points {
		dx := point.x - centroids[0].x
		dy := point.y - centroids[0].y
		minD := dx*dx + dy*dy

		for _, center := range centroids {
			dx := point.x - center.x
			dy := point.y - center.y
			d := dx*dx + dy*dy
			if d < minD {
				minD = d
			}
		}

		D[idx].D2 = minD
		D[idx].pointIndex = idx
	}
}

type distSlice []dist

func (ps distSlice) Len() int           { return len(ps) }
func (ps distSlice) Less(i, j int) bool { return ps[i].D2 < ps[j].D2 }
func (ps distSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
