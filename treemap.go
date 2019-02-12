/*
Package treemap presents a way to visualize hierarchical information
structures efficiently in 2-D display surface. treemap works with
types that are defined recursively and have weight/size defined
for each node. If a node implements the TreeMaper interface
    type TreeMaper interface {
        Identity() string
        Weight() float64
        Descendants() []TreeMaper
    }
then it can be drawn
This weight may represent a single domain property
(such as disk usage or file age for directory tree)
A node's weight will determine its display size and only the
space represented by leaf nodes are shown.

Properties:

If node i is an ancestor of node j, then the bounding box of node i
completely encloses or is equal to the bounding box of node j.

The bounding boxes of two nodes intersect iff one node is an
ancestor of the other

Nodes occupy a display area proportional to their weight.

The weight of a node is greater than or equal to the sum of its
children.

Algorithm:

The procedure tracks the cursor movement in the tree

1. The root node draws itself within its rectangular bounds
2. Checks that maximum depth to draw is not reached.
2. Sets new bounds and drawing properties for each of its children,
inverting the orientation for the next drawing phase.
3. Each of its children do steps 1, 2 and 3 recursivley till
the complete tree is mapped.

See the paper published by the original authors for details
https://www.cs.umd.edu/~ben/papers/Johnson1991Tree.pdf
*/
package treemap

import (
	"fmt"
	"image"
	"io"

	svg "github.com/ajstarks/svgo"
)

// Orientation is the diretion to start a slice
type Orientation string

const (
	// Horizontal describes a line along the x-axis
	Horizontal Orientation = "horizontal"

	// Vertical describes a line along the y-axis
	Vertical = "vertical"
)

// TreeMapper interface represents a hierarchical
// data structure that can be drawn to form a treemap
type TreeMapper interface {
	Identity() string
	Weight() float64
	Descendants() []TreeMapper
}

func drawTree(
	t TreeMapper,
	svg *svg.SVG,
	path Orientation,
	bound image.Rectangle,
	depth, maxDepth int,
) {
	// check that maxDepth is not reach
	if maxDepth != 0 && depth >= maxDepth {
		return
	}
	// consumed is the unit of width or height consumed
	var consumed float64
	parentWeight := t.Weight()
	nextPath := Horizontal
	if path == Horizontal {
		nextPath = Vertical
	}

	for _, c := range t.Descendants() {
		var proportion float64
		var newBound image.Rectangle
		var color string
		if path == Horizontal {
			// slicing would be along y-axis
			// x values may not be touched ?
			// `proportion` to consume is c.Weight() / parentWeight
			// `consumed` will  determine the offset to start new consumption
			// `proportion` tells the unit of width or height to consume
			// `bound` is (x0, y0) - (x1, y1)
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
			proportion = (c.Weight() / parentWeight) * float64(bound.Dy())
			x0 := bound.Min.X
			x1 := bound.Max.X
			y0 := bound.Min.Y + int(consumed+0.5)
			y1 := bound.Min.Y + int(consumed+0.5) + int(proportion+0.5)
			min := image.Point{x0, y0}
			max := image.Point{x1, y1}
			newBound = image.Rectangle{min, max}
		} else {
			// slicing would be along the y-axis
			// x0 -> parentX0 + consumed
			// x1 -> parentX0 + consumed + proportion
			// y0 -> parentY0
			// y1 -> parentY1
			proportion = (c.Weight() / parentWeight) * float64(bound.Dx())
			x0 := bound.Min.X + int(consumed+0.5)
			x1 := bound.Min.X + int(consumed+0.5) + int(proportion+0.5)
			y0 := bound.Min.Y
			y1 := bound.Max.Y
			min := image.Point{x0, y0}
			max := image.Point{x1, y1}
			newBound = image.Rectangle{min, max}
		}
		color = newRgb(
			int(parentWeight)>>uint(2),
			int(parentWeight)>>uint(1),
			int(parentWeight+proportion),
		).String()

		drawNode(
			svg,
			c.Identity(),
			newBound,
			color,
		)

		// update consumed for the next iteration
		// then send child to draw itself
		consumed += proportion
		drawTree(c, svg, nextPath, newBound, depth+1, maxDepth)
	}
}

// drawNode draws a treemap node using the bound,
// color, an identity passed in to create an svg element
func drawNode(
	svg *svg.SVG,
	identity string,
	bound image.Rectangle,
	color string,
) {
	svg.Rect(
		bound.Min.X,
		bound.Min.Y,
		bound.Dx(),
		bound.Dy(),
		"fill: "+color+";stroke: #fff;",
	)
	svg.Text(
		bound.Min.X,
		bound.Min.Y+10,
		identity,
		"font-size:10px;padding:30px;text-anchor: start;",
	)
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

// rgb implements Stringer interface and returns
// the svg color notation for a node in the form
// rgb(#, #, #) where # is the corresponding componenent value
func (c rgb) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.r, c.g, c.b)
}

// DrawTreemap draws the tree-map described by treemaper
// and writes the resulting tree-map to the io.Writer
// at a depth less than or equal to the maxDepth
// from the given start orientation
func DrawTreemap(
	w io.Writer,
	tm TreeMapper,
	width, height int,
	startPath Orientation,
	maxDepth int,
) {
	svg := svg.New(w)
	svg.Start(width, height)
	bound := image.Rect(0, 0, width, height)
	drawTree(tm, svg, startPath, bound, 0, maxDepth)
	svg.End()
}
