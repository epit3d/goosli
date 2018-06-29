# goosli
Package to implement your own slicer using existing possibilities: cutting, crossing, 
slicing by vector, mesh simplification, line simplification etc.

For example you can slice by profile (for <b>5axes 3d printer</b>) - find center of 
each layer, make line from points, simplify line - and slice by this line. 
See slicers/slice_by_profile.go. For sure your 3d-printer have to support bed rotations.

Feel free to open issues or implement your slicing algorithms.

Do not forget to place <b>data directory</b> near your binary. 

Here is a viewer to see results of slicing - [spycer](https://github.com/l1va/spycer).

# goosli-cutter
Cut stl in two stls by required plane.

# goosli-simplifier
Simplifies stl to required count of triangles.

### Thanks
A lot of ideas and code was taken from various [fogleman](https://github.com/fogleman) 
repos. Thank you!
