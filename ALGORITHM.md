Generation Algorithm
===

The infinite rail world is split into tiles, each of which represents 100m. The
general idea of the algorithm `generate` is to be able to generate any tile of the world (identified by its number)
without generating the whole (infinite) world before. Thus, `generate` takes two parameters: a seed and the tile number.

The output of `generate(seed, tile)` has to meet the following two requirements:
1. As long as `seed` and `tile` are identical, all calls to `generate` must return the same output.
2. Given the results of `generate(seed, tile)` and `generate(seed, tile + 1)`, both results should be compatible to each other.

The implementation of `generate` can be roughly divided into the following steps. They are
explained in detail in the following sections.

1. Determine the segment of `tile`. A _segment_ is a structure that spans over more than on tile, e.g. a train station.
2. Generate the structure inside the segment.
3. If `tile` is not the rightmost tile in the segment, return the tile. Else, proceed to 4.
4. Modify the tile such that it is compatible to the first tile of the next segment, by generating the next segment.

Segments
---

A segment is a rail structure, that spans over more than one tile. Every tile belongs to one segment. Examples for segments
are stations, straight track, rail yards etc.

To compute the segment borders, the theoretical meaning of "infinite" is altered to "so much that it is practically infinite".
For example, a number of 30,000 tiles of 100m equals a quarter of the earth's perimeter. Let `sampling` be the maximal number
of tiles.

Given the seed, a `list` of random points in the interval `[0,1)` is generated. Let `x` be the 
result of `tile * sampling / len(list)`, where `tile` is the number of the requested tile. Let `min = list[floor(x)]` and `max = list[ceil(x)]`.

The values between `min` and `max` are interpolated using the Smoothstep function and its result multiplied by 10. This function `f` is used to assign the tile
a number between 0 and 10.

The above computations have the advantage, that `f(tile-1)`, `f(tile)`, and `f(tile + 1)` are different, yet close by.
To find the start of the segment in which `tile` resides, it suffices to find a `start` such that `ceil(f(tile - start)) - ceil(f(tile)) != 0`.
In the same way, the end of the segment can be found.

Thus, for `tile`, the start and the length of the segment is known. A pseudorandom generator `random`
is initialized with `f(start-index + 1)`. The generator determines how the segment will look like.


Segment Generation
---

At this point, the length of the tile's segment is known. Now, the type of the segment
is specified by utilizing `random` and taking the segment's length into account. The following
types are supported:

* straight track
* junctions
* stations

Each of these type comes with a probability and, possibly, a length restriction. For example,
a station may not be generated if the segment is 100 tiles long.

The segment is then instantiated using a localized algorithm specified for the segment type. Then,
the requested type within the segment is returned. If the tile was not the last tile of the segment, it
can be rendered directly, otherwise, it is made compatible to the next segment.


