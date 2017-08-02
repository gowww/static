# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) static [![GoDoc](https://godoc.org/github.com/gowww/static?status.svg)](https://godoc.org/github.com/gowww/static) [![Build](https://travis-ci.org/gowww/static.svg?branch=master)](https://travis-ci.org/gowww/static) [![Coverage](https://coveralls.io/repos/github/gowww/static/badge.svg?branch=master)](https://coveralls.io/github/gowww/static?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/static)](https://goreportcard.com/report/github.com/gowww/static) ![Status Stable](https://img.shields.io/badge/status-stable-brightgreen.svg)

Package [static](https://godoc.org/github.com/gowww/static) provides a handler for static file serving with cache control and automatic fingerprinting.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/static
	```

2. Import it in your code:

	```Go
	import "github.com/gowww/static"
	```

## Usage

Use [Handle](https://godoc.org/github.com/gowww/static#Handle) with the URL path prefix and the source directory to get a [Handler](https://godoc.org/github.com/gowww/static#Handler) that will serve your static files:

```Go
staticHandler := static.Handle("/static/", "static")

http.Handle("/static/", staticHandler)
```

Use [Handler.Hash](https://godoc.org/github.com/gowww/static#Handler.Hash) to append the file fingerprint to a file name (if the file can be opened, obviously):

```Go
staticHandler.Hash("scripts/main.js")
```

But generally, you'd want to use this method in your templates:

```Go
tmpl := `<script src="{{asset "scripts/main.js"}}"></script>`

views := template.Must(template.New("main").Funcs(template.FuncMap{
	"asset": staticHandler.Hash,
}).Parse(tmpl))
```

## References

- [Strategies for cache-busting CSS — CSS Tricks](https://css-tricks.com/strategies-for-cache-busting-css/)
- [Fingerprinting images to improve page load speed — Imgix](https://docs.imgix.com/best-practices/fingerprinting-images-improve-page-load-speed)
- [Revving filenames: don’t use querystring](http://www.stevesouders.com/blog/2008/08/23/revving-filenames-dont-use-querystring/)
