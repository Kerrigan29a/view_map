# View CSV

This an example program to plot a world map using [drawille-go](https://github.com/Kerrigan29a/drawille-go).
To render the world map this program uses [world-50m.txt](world-50m.txt) and a [Cylindrical Equidistant Projection](http://mathworld.wolfram.com/CylindricalEquidistantProjection.html)

# Documentation

Documentation is available at [godoc](https://godoc.org/github.com/Kerrigan29a/view_map)

# Example

```shell
make run_flat SCALE=1
make run_miller37 SCALE=1
make run_miller43 SCALE=1
make run_miller50 SCALE=1

make run ARGS="-w world-50m.txt -p flat -s 0.5"
```
