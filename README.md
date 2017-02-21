# tspos-lbtw: Liquor Bottle Tare Weight microservice

[![Build Status](http://img.shields.io/travis/wormling/tspos-lbtw.svg?branch=master)](https://travis-ci.org/wormling/tspos-lbtw) [![GO Report Card](https://goreportcard.com/badge/github.com/wormling/tspos-lbtw)](https://goreportcard.com/report/github.com/wormling/tspos-lbtw)

REST API for liquor inventory using tare weights.

Requirements:
* [Go](https://golang.org/)
* [MongoDB 3.4+](https://www.mongodb.com/)

Dependencies:
* [Gin Gonic](https://github.com/gin-gonic) 
* [mgo](https://labix.org/mgo)

To install, run:
```bash
$ go get github.com/wormling/tspos-lbtw.v1
```

To configure edit config.yaml:
```yaml
core:
  listener:
    bind: localhost
    port: 2727
  database:
    url: mongodb://localhost:27017/tspos_lbtw
```
