package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/willpoint/treemap"
)

// TNode is a treemap node
type TNode struct {
	Name     string   `json:"name"`
	Size     float64  `json:"size,omitempty"`
	Children []*TNode `json:"children,omitempty"`
}

var _ treemap.TreeMapper = &TNode{}

// Identity implements the TreeMapper interface for TNode
func (t *TNode) Identity() string {
	return t.Name
}

// Weight implements the TreeMapper interface for TNode
func (t *TNode) Weight() float64 {
	var sum float64
	var weight func(t *TNode)
	weight = func(t *TNode) {
		if t.Children != nil {
			for _, d := range t.Children {
				weight(d)
			}
		} else {
			sum += t.Size
		}
	}
	weight(t)
	return sum
}

// Descendants implements the TreeMaper interface for TNode
func (t *TNode) Descendants() []treemap.TreeMapper {
	retv := []treemap.TreeMapper{}
	for _, i := range t.Children {
		retv = append(retv, i)
	}
	return retv
}

func main() {

	width := flag.Int("w", 800, "width of rectange")
	height := flag.Int("h", 600, "height of rectangle")
	infile := flag.String("in", "", "filename to get data (json file)")
	outfile := flag.String("out", "output.svg", "filename to save data (in svg)")
	maxDepth := flag.Int("depth", 0, "max depth to draw the treemap")
	flag.Parse()

	if *infile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// output to save visualization
	out, err := os.OpenFile(*outfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("opening file", err)
	}

	tmap := new(TNode)

	// data to visualize
	f, err := os.Open(*infile)
	if err != nil {
		log.Fatal("opening file: ", err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(tmap)

	var orientation treemap.Orientation
	orientation = treemap.Vertical
	if *width < *height {
		orientation = treemap.Horizontal
	}
	treemap.DrawTreemap(out, tmap, *width, *height, orientation, *maxDepth)

}
