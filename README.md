#### GIN dumper

[![codecov](https://codecov.io/gh/nikonm/gin-dumper/branch/master/graph/badge.svg)](https://codecov.io/gh/nikonm/gin-dumper)
[![Build Status](https://travis-ci.org/nikonm/gin-dumper.svg?branch=master)](https://travis-ci.org/nikonm/gin-dumper)

 - Gin middleware handler for dumping body and headers of request(and response). 
 - You can easily use it with your logger

###### Sample
 
````
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/nikonm/gin-dumper"
    "fmt"
)

func main () {
    router := gin.New()

    opt := &gin_dumper.Options{TrimNewLineInRequest: "   "}

    router.Use(gin_dumper.Dumper(func(output *gin_dumper.Output) {
        fmt.Printf("REQUEST: %s %s Headers: %v Body: %s Resp: %d JsonBody: %s\n",
			output.Request.Method,
			output.Request.Url.RequestURI(),
			output.Request.Headers,
			output.Request.Body,
			output.Response.StatusCode,
			output.Response.Body,
		)
    }, opt))

    router.POST("/test", func(c *gin.Context) {
		// Some Action
    })
	
    router.Run(":80")
}

````
