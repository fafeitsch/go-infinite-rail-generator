The Generation Algorithm
===
This file explains the overall generation algorithm:

__Given__ a tile number _n_, and a seed _s_, generate a tile that:
* matches its left neighbour _n – 1_ in terms of tracks, stations, etc.
* matches its right neighbour _n + 1_ in terms of tracks, stations, etc.
* looks identical as long as the seed _s_ is not changed. For different seeds _s_ the tile should
    look different

The last point is no hard requirement because there a lot more seeds than possible tile configurations.

Phase 1
----

For any given seed _s_, 512 points in the range `[0,1]` are generated. This
is called the _noise_. The noise can be regarded as a function `[0, 511] → [0,1]`. This, however,
only allows different tiles, and every tile could differ vastly from its neighbouring tiles.

We define 30000 as the arbitrary limit of different tiles per seed. Then, the range `[0,511]`
is split into 30000 equal units. For any tile number _n < 30000_, we can now compute both points
between which the tile lies. The number 30000 is called _sampling_.

The value of both points are interpolated with the [Smoothstep](https://en.wikipedia.org/wiki/Smoothstep)
function and this function is used to get a raw fractional value for the requested tile _n_. The advantage
of this approach is that neighbouring tiles _n - 1_ and _n + 1_ are having similar, but slightly
different values.

Phase 2
---

Now, we can start building our tile. First, we determine the number of tracks by tacking the
fractional value, multiplying it with some constants and rounding to the nearest integer.

Then we determine whether the tile contains a station. Since we want stations to span several tiles
we also use the algorithm from the first phase, but with a slightly derived seed and smaller sampling.
If the fractional value of the tile is below a certain threshold, then its a station tile. Again, if a tile
is a station, the probability that the neighbouring tiles are stations, too, is heigh.

