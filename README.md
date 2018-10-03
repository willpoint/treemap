# Tree-Map

An implementation of Tree-Map in Go for visualizing hierarchical data.

An example usage with a compiled binary is
```
`treemap -in=data.json -depth=2 -out=output.svg
```

**Todo**
- Add an API that explicitly tells treemap to `DrawToSVG(/** an svg interface */)` to enable user easily swap out the svg implementation used internally.
- Add a method to draw to a pre-existing image.
