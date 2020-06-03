# MyHTTP

MyHTTP makes http requests and prints the address of the request along with
the MD5 hash of the response.

## Requirements

- Go 1.10 or more recent
- Make (*optional*)

## Installing

```
go get github.com/scorphus/myhttp
```

Once installed, you'll find `myhttp` on your `$PATH`, assuming you [installed
and configured Go accordingly](https://golang.org/doc/install).

## Building

[Makefile](Makefile) offers a handy `build` target:

```
make build
```

This target builds the application with the following command:

```
go build -x -v -race
```

## Running

Once built, you can find the `myhttp` executable in the current directory.

```
$ ./myhttp --help
Usage: myhttp [flags] [url ...]
  -parallel int
    	limit the number of parallel requests (default 10)
```

### Exempli gratia

```
$ ./myhttp globo.com dw.com debian.org
http://debian.org ab590dac91d0850169462c9fad24c4d3
http://dw.com dfeb459761932def97ea0bf633bba56a
http://globo.com a4148052ca14220a8637622d5c5aa7a3
```

Using `parallel`:

```
$ ./myhttp -parallel 2 globo.com dw.com debian.org
http://dw.com dfeb459761932def97ea0bf633bba56a
http://globo.com a4148052ca14220a8637622d5c5aa7a3
http://debian.org ab590dac91d0850169462c9fad24c4d3
```

## Testing

[Makefile](Makefile) also has targets to run tests:

```
make test
```

To run tests and display test coverage on the console:

```
make cover
```

To run tests and display test coverage on the console and on a browser:

```
make cover-html
```

> *p.s.: run `make list` to list all available targets*

## License

Code in this repository is distributed under the terms of the 3-Clause BSD
License (BSD-3-Clause).

See [LICENSE](LICENSE) for details.
