package dbscan

import (
	"math"
	"strconv"

	"github.com/kyroy/kdtree"
)

type Cluster struct {
	kdt       *kdtree.KDTree
	minRange  float64
	miniPoint int
}

func NewCluster(p []kdtree.Point, minRange float64, miniPoint int) *Cluster {

	return &Cluster{
		kdt:       kdtree.New(p),
		minRange:  minRange,
		miniPoint: miniPoint,
	}
}

func (c *Cluster) Get() (result [][]kdtree.Point) {
	c.kdt.Balance()
	result = make([][]kdtree.Point, 0)
	visited := map[string]struct{}{}
	for _, point := range c.kdt.Points() {
		pointId := GetId(point)
		_, isVisited := visited[pointId]
		if isVisited {
			continue
		}

		neighbours := make([]kdtree.Point, 0)
		c.findNeighbours(point, visited, &neighbours)
		if len(neighbours)+1 >= int(c.miniPoint) {
			visited[pointId] = struct{}{}
			neighbours = append(neighbours, point)
			result = append(result, neighbours)
		}
	}
	return result
}

func (c *Cluster) findNeighbours(point kdtree.Point, visited map[string]struct{}, ret *[]kdtree.Point) {
	neighbours := make([]kdtree.Point, 0)
	neighboursPoints := c.kdt.KNN(point, int(c.miniPoint))
	for _, potNeigb := range neighboursPoints {
		if GetId(point) != GetId(potNeigb) && distance(point, potNeigb) <= c.minRange {
			neighbours = append(neighbours, potNeigb)
		}
	}
	if len(neighbours) == c.miniPoint-1 {
		visited[GetId(point)] = struct{}{}
		for _, npoint := range neighbours {
			pointId := GetId(npoint)
			_, isVisited := visited[pointId]
			if isVisited {
				continue
			}
			visited[GetId(point)] = struct{}{}
			*ret = append(*ret, npoint)
			c.findNeighbours(npoint, visited, ret)
		}
	}
}

func GetId(p1 kdtree.Point) string {
	ret := ""
	for i := range make([]struct{}, p1.Dimensions()) {
		ret += strconv.FormatFloat(p1.Dimension(i), 'e', 16, 64) + ","
	}
	return ret
}

func distance(p1, p2 kdtree.Point) float64 {
	sum := 0.
	for i := 0; i < p1.Dimensions(); i++ {
		sum += math.Pow(p1.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}
