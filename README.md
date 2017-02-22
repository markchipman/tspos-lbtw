# tspos-lbtw: Liquor Bottle Tare Weight microservice

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](http://img.shields.io/travis/wormling/tspos-lbtw.svg?branch=master)](https://travis-ci.org/wormling/tspos-lbtw) 
[![GO Report Card](https://goreportcard.com/badge/github.com/wormling/tspos-lbtw)](https://goreportcard.com/report/github.com/wormling/tspos-lbtw)
[![Coverage Status](https://coveralls.io/repos/wormling/tspos-lbtw/badge.png?branch=v1)](https://coveralls.io/r/wormling/tspos-lbtw)

REST API for liquor inventory using tare weights.

Requirements:
* [Go](https://golang.org/)
* [MongoDB 3.4+](https://www.mongodb.com/)

Dependencies:
* [Gin Gonic](https://github.com/gin-gonic) 
* [mgo](https://labix.org/mgo)

Package Management:
* [Glide](https://github.com/Masterminds/glide)

Testing:
* [Ginkgo](https://onsi.github.io/ginkgo/)
* [Gomega](http://onsi.github.io/gomega/)

To configure edit config.yaml:
```yaml
core:
  listener:
    bind: localhost
    port: 2727
  database:
    url: mongodb://localhost:27017/tspos_lbtw
```
