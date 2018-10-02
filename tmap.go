package main

import (
	"encoding/json"
	"flag"
	"image"
	"log"
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// A TreeMap presents hierarchical information structures
// efficiently in 2-D display surface.
// TM requires that a weight be assigned to each node
// this weight may represent a single domain property
// (such as disk usage or file age for a directory tree)
// To simplify this implementation size is used to describe
// this weight.
// A node's weight (bounding box) will determine its display
// size and can be thought of as a measure of importance
// or degree of interest.
// The following are properties of a treemap
//
// > If node i is an ancestor of node j, then the bounding box
// of node i completely encloses or is equal to the bounding box
// of j
// > the bouding boxes of two nodes intersect iff one node
// is an ancestor of the other
// > nodes occupy a display area proportional to their weight.
// > The weight of a node is greater than or equal to the sum
// > of the weight of its children
//
//  Displaying visualization
// Once the bounding box of a node is set, a variety of display
// props determine how the node is drawn within it.
// > color (hue, saturation, brightness)
// > texture, shape, border, blinking, etc...
//
// Algorithm
// The procedure will draw a Tree-Map and track the cursor
// movement in the tree.
// >
// 1. The node draws itself within it's rectangular bounds
// according to its display property (weight, color, border)
// 2. The node sets new bounds and drawing properties for each of
// its children, and recursively sends each child a drawing command
// The bounds of a node's children form either a vertical or horizontal
// partitioning of the display space allocated to the node
// >
//

const (
	// HORIZONTAL CONSTANT
	HORIZONTAL = "horizontal"
	// VERTICAL CONSTANT
	VERTICAL = "vertical"
)

// todo(uz)
// create an interface type to describe elements that
// can be drawn to form a treemap - eg. Size() string | Name() string

// TNode is a treemap node
type TNode struct {
	Name     string   `json:"name"`
	Size     float64  `json:"size,omitempty"`
	Children []*TNode `json:"children,omitempty"`

	color       string
	orientation string
	depth       int
	bound       image.Rectangle
}

// drawNode uses information sent from
// parent to correctly draw itself
func (t *TNode) drawNode(
	svg *svg.SVG,
	bound image.Rectangle,
	orientation string,
	color string,
	depth int,
) {
	t.depth = depth
	t.orientation = orientation
	t.color = color
	t.bound = bound
	svg.Rect(
		bound.Min.X,
		bound.Min.Y,
		bound.Dx(),
		bound.Dy(),
		"fill: "+t.color+";stroke: #fff;",
	)
	svg.Text(
		bound.Min.X,
		bound.Min.Y+10,
		t.Name,
		"font-size:10px;padding:30px;text-anchor: start;",
	)
}

func (t *TNode) size() float64 {
	var sum float64
	each([]*TNode{t}, func(n *TNode) {
		sum += n.Size
	}, nil)
	return sum
}

func (t *TNode) drawTree(svg *svg.SVG, maxDepth int) {

	// check that maxDepth is not reach
	if maxDepth != 0 && t.depth >= maxDepth {
		return
	}

	// consumed is the unit of width or height consumed
	var consumed float64
	mSize := t.size()
	var nextOrientation string
	if t.orientation == VERTICAL {
		nextOrientation = HORIZONTAL
	} else {
		nextOrientation = VERTICAL
	}
	// create rectangular bound for each child
	for _, c := range t.Children {
		var proportion float64
		var bound image.Rectangle
		var color string
		if t.orientation == HORIZONTAL {
			// slicing would be along y-axis
			// x values may not be touched ?
			// proportion to consume is c.size / mSize
			// `consumed` will tell the determine the offset
			// to start new consumption
			// `proportion` tells the unit of width or height to
			// consume
			// bound is (x0, y0)-(x1, y1)
			//
			// (x0, y0)				(x1, y0)
			//    +--------------------+
			//	  |					   |
			//	  |					   |
			//	  |					   |
			//	  |					   |
			//	  |					   |
			// 	  +--------------------+
			//	(x0, y1)			(x1, y1)
			//
			// set values for all points in the rect
			// x0 -> parentX0
			// x1 -> parentX1
			// y0 -> parentY0 + consumed
			// y1 -> parentY0 + consumed + proportion
			proportion = (c.size() / mSize) * float64(t.bound.Dy())
			x0 := t.bound.Min.X
			x1 := t.bound.Max.X
			y0 := t.bound.Min.Y + int(consumed+0.5)
			y1 := t.bound.Min.Y + int(consumed+0.5) + int(proportion+0.5)
			min := image.Point{x0, y0}
			max := image.Point{x1, y1}
			bound = image.Rectangle{min, max}
		} else {
			// slicing would be along the y-axis
			// x0 -> parentX0 + consumed
			// x1 -> parentX0 + consumed + proportion
			// y0 -> parentY0
			// y1 -> parentY1
			proportion = (c.size() / mSize) * float64(t.bound.Dx())
			x0 := t.bound.Min.X + int(consumed+0.5)
			x1 := t.bound.Min.X + int(consumed+0.5) + int(proportion+0.5)
			y0 := t.bound.Min.Y
			y1 := t.bound.Max.Y
			min := image.Point{x0, y0}
			max := image.Point{x1, y1}
			bound = image.Rectangle{min, max}
		}
		color = newRgb(
			int(mSize)>>uint(2),
			int(mSize)>>uint(1),
			int(mSize+proportion),
		).String()
		c.drawNode(
			svg,
			bound,
			nextOrientation,
			color,
			t.depth+1,
		)
		// update consumed for the next iteration
		// then send child to draw itself
		consumed += proportion
		c.drawTree(svg, maxDepth)
	}
}

// helper to run before and after func for each node in the tree
func each(nn []*TNode, before, after func(t *TNode)) {
	for _, c := range nn {
		if before != nil {
			before(c)
		}
		if c.Children != nil {
			each(c.Children, before, after)
		}
		if after != nil {
			after(c)
		}
	}
}

// rgb is the color model used for the treemap
type rgb struct {
	r, g, b uint8
}

// newRgb create a new color model
func newRgb(r, g, b int) rgb {
	return rgb{
		r: uint8(r & 0xff),
		g: uint8(g & 0xff),
		b: uint8(b & 0xff),
	}
}

// itoa helps convert an rgb component value to string
func itoa(n uint8) string {
	return strconv.Itoa(int(n))
}

// rgb implements Stringer interface and returns
// the svg color notation for a node in the form
// rgb(#, #, #) where # is the corresponding componenent value
func (c rgb) String() string {
	return "rgb(" + itoa(c.r) + "," +
		itoa(c.g) + "," + itoa(c.b) + ")"
}

func main() {

	// commandline arguments
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
	svg := svg.New(out)
	svg.Start(*width, *height)
	tmap := new(TNode)

	// data to visualize
	f, err := os.Open(*infile)
	if err != nil {
		log.Fatal("opening file: ", err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(tmap)
	rect := image.Rect(0, 0, *width, *height)
	tmap.orientation = VERTICAL
	tmap.bound = rect
	if *width < *height {
		tmap.orientation = HORIZONTAL
	}
	tmap.drawTree(svg, *maxDepth)
	svg.End()
}
