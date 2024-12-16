package dbscan

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func GenPoints() []kdtree.Point {
	maxLen := 2000
	ret := make([]kdtree.Point, 0)
	var dd map[string]struct{} = make(map[string]struct{})
	for _ = range make([]struct{}, maxLen) {
		p := &points.Point2D{
			X: float64(rand.Intn(200)),
			Y: float64(rand.Intn(200)),
		}
		if _, ok := dd[GetId(p)]; !ok {
			dd[GetId(p)] = struct{}{}
			ret = append(ret, p)
		}
	}
	return ret
}

func ParsePoints() []kdtree.Point {
	ret := make([]kdtree.Point, 0)
	buf, _ := os.ReadFile("testdata")
	data := strings.Split(string(buf), "\n")
	xdata := strings.Split(data[0], "\t")
	ydata := strings.Split(data[1], "\t")
	for i := range xdata {
		xx, _ := strconv.ParseFloat(xdata[i], 64)
		yy, _ := strconv.ParseFloat(ydata[i], 64)
		p := &points.Point2D{
			X: xx,
			Y: yy,
		}
		ret = append(ret, p)
	}
	return ret
}

func TestGenPic(t *testing.T) {
	pp := plot.New()
	pp.Title.Text = "data"
	pp.X.Label.Text = "X"
	pp.Y.Label.Text = "Y"
	r := NewCluster(ParsePoints(), 0.5, 10).Get()
	data := make([]interface{}, 0)
	for _, v := range r {
		sdata := make(plotter.XYs, 0)
		for _, point := range v {
			sdata = append(sdata, plotter.XY{
				X: point.Dimension(0),
				Y: point.Dimension(1),
			})
		}
		if len(sdata) > 0 {
			data = append(data, sdata)
		}
	}
	plotutil.AddScatters(pp, data...)
	if err := pp.Save(4*vg.Inch, 4*vg.Inch, "test.png"); err != nil {
		t.Fatal(err)
	}

}
