# DataHose

# Installation

```
$ go get -u github.com/bluele/datahose
```

# Example

```
$ go run examples/app.go
2015/05/17 13:25:34 serving at :8000

# open other terminal
$ curl http://localhost:8000
{"time":"2015-05-17 13:20:17.601006083 +0900 JST"}
{"time":"2015-05-17 13:20:17.601072586 +0900 JST"}
{"time":"2015-05-17 13:20:17.601077136 +0900 JST"}
{"time":"2015-05-17 13:20:17.601080451 +0900 JST"}
{"time":"2015-05-17 13:20:17.60108329 +0900 JST"}
{"time":"2015-05-17 13:20:17.60108619 +0900 JST"}
{"time":"2015-05-17 13:20:17.60108892 +0900 JST"}
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>