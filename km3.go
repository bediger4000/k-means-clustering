package main

/*
   K-means clustering with k-means++ initial centroid choice.
*/

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// Point an x,y cartesian point
type Point struct {
	pop float64
	x   float64
	y   float64
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

	for i := range points {
		fmt.Printf("%4d  %.0f  %f %f\n", i, points[i].pop, points[i].x, points[i].y)
	}

	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

	kmeanscluster(k, points)
}

func kmeanscluster(k int, points []Point) {

	/*
	   Initialization:
	   1. Compute the desired cluster size, n/k.
	   2. Initialize means, preferably with k-means++
	   3. Order points by the distance to their nearest cluster minus distance to
	      the farthest cluster (= biggest benefit of best over worst assignment)
	   4. Assign points to their preferred cluster until this cluster is full, then
	      resort remaining objects, without taking the full cluster into account
	      anymore

	*/

	// clusterSize := len(points) / k

	centroids := kMeansPPCentroids(k, points)
	fmt.Printf("Final centroids %d:\n", len(centroids))
	for i := range centroids {
		fmt.Printf("centroid %d (%f, %f)\n", i, centroids[i].x, centroids[i].y)
	}

	orderedPoints := orderPointsByDistance(points, centroids)
	fmt.Printf("%d orderedPoints\n", len(orderedPoints))
	for i := range orderedPoints {
		pt := points[orderedPoints[i].pointIndex]

		fmt.Printf("  ordered point %d, dist diff %f, point %d (%f, %f), %v\n",
			i,
			orderedPoints[i].distDiff,
			orderedPoints[i].pointIndex,
			pt.x, pt.y, orderedPoints[i].distances,
		)
	}

	assignToClusters(orderedPoints, k)

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
		n, err := fmt.Fscanf(fin, "%f %f %f\n", &p.pop, &p.x, &p.y)
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

   I'm not pretending this is optimal - and it may not even be correct,
   but it's the only way I could think of to get a probability proptional
   to D^2

   Parameter D already has dx^2+dy^2 as the D2 element value.
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
			fmt.Printf("Sum of D62 %f, random value %f, picking interval %d/%f\n",
				sumValues, inInterval, i, intervals[i].maxval)
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

	fmt.Printf("fill distances, %d Points, %d dist, %d centroids\n", len(points), len(D), len(centroids))
	fmt.Printf("Centroids:\n")
	for i := range centroids {
		fmt.Printf("Centroid %d (%f, %f)\n", i, centroids[i].x, centroids[i].y)
	}
	for i := range D {
		fmt.Printf("dist %d: Point %d (%f, %f) dist %f\n", i, D[i].pointIndex,
			points[D[i].pointIndex].x, points[D[i].pointIndex].y, D[i].D2)
	}
}

type pointDist struct {
	distDiff   float64
	distances  []dist // here, the pointIndex is that of a centroid
	pointIndex int
}

func assignToClusters(orderedPoints []pointDist, k int) (clusters [][]Point) {
	return
}

/*
Order points by the distance to their nearest cluster minus distance to
the farthest cluster (= biggest benefit of best over worst assignment)
*/
func orderPointsByDistance(points []Point, centroids []Point) []pointDist {
	var orderedPoints []pointDist
	for i := range points {
		oPoint := pointDist{pointIndex: i}
		for j := range centroids {
			dx := points[i].x - centroids[j].x
			dy := points[i].y - centroids[j].y
			distToCentroid := math.Sqrt(dx*dx + dy*dy)
			oPoint.distances = append(oPoint.distances, dist{D2: distToCentroid, pointIndex: j})
		}
		ds := distSlice(oPoint.distances)
		sort.Sort(ds)
		m := len(centroids)
		oPoint.distDiff = oPoint.distances[m-1].D2 - oPoint.distances[0].D2

		orderedPoints = append(orderedPoints, oPoint)
	}

	sort.Sort(pointDistSlice(orderedPoints))

	return orderedPoints
}

type distSlice []dist

func (ps distSlice) Len() int           { return len(ps) }
func (ps distSlice) Less(i, j int) bool { return ps[i].D2 < ps[j].D2 }
func (ps distSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

type pointDistSlice []pointDist

func (pds pointDistSlice) Len() int           { return len(pds) }
func (pds pointDistSlice) Less(i, j int) bool { return pds[i].distDiff < pds[j].distDiff }
func (pds pointDistSlice) Swap(i, j int)      { pds[i], pds[j] = pds[j], pds[i] }
