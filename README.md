go-extra-lib
============

go-extra-lib is a package with Go construct commonly used in a lot of my projects.
To avoid copy-pasting I have gathered them in a single importable package.


Packages
--------

- `slices2`: Extension to [golang.org/x/exp/slices](https://pkg.go.dev/golang.org/x/exp/slices), with more slice tricks.
- `testing2`: Various assertions for writing tests, automatically-generated
    - `assert`: assertions which immediately fail a test
    - `check`: assertions which allow a test to continue

License
-------

go-extra-lib is provided under the MIT license, included in the <LICENSE> file.
