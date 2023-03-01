go-extra-lib
============

go-extra-lib is a package with Go construct commonly used in a lot of my projects.
To avoid copy-pasting I have gathered them in a single importable package.


Packages
--------

- `container`:
    - `bitset`: An efficient implementation of a set of unsigned numbers
- `encoding`:
    - `mcsv`: CSV, but map\[string\]string instead of \[\]string
- `io2`: Occasionally useful io.Readers
- `iter`: Generic iterators and operations on such iterators
- `slices2`: Extension to [golang.org/x/exp/slices](https://pkg.go.dev/golang.org/x/exp/slices), with more slice tricks.
- `testing2`: Various assertions for writing tests, automatically-generated
    - `assert`: assertions which immediately fail a test
    - `check`: assertions which allow a test to continue

TODO
----

- [ ] `container/gheap`: Generic version of `container/heap`
- [ ] `container/glist`: Generic version of `container/list`
- [ ] `container/gring`: Generic version of `container/ring`
- [ ] `container/set`: An unordered collection of elements (map\[T\]struct{})
- [ ] `iter/stream`: Java Stream-like wrapper on iterator operations
- [ ] `maps2`: Extension to [golang.org/x/exp/maps](https://pkg.go.dev/golang.org/x/exp/maps), with more operations on maps.
- [ ] `matrix`: 2D matrices of numbers

License
-------

go-extra-lib is provided under the MIT license, included in the <LICENSE> file.
