/*
Package treemap presents a way to visualize hierarchical information
structures efficiently in 2-D display surface. treemap works with
types that are defined recursively and have weight/size defined
for each node. If a node implements the TreeMaper interface
    type TreeMaper interface {
        Name() string
        Weight() float64
        Children() []TreeMaper
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

See the paper published by the original author for details
https://www.cs.umd.edu/~ben/papers/Johnson1991Tree.pdf
*/
package treemap
