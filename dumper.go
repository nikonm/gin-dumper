package gin_dumper

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	Url     *url.URL
	Headers http.Header
	Body    string
}

type Response struct {
	StatusCode int
	Body       string
}

type Output struct {
	Request  *Request
	Response *Response
}

type Options struct {
	TrimNewLineInRequest string // Replace "\n" in request body to this value, if set. Not worked in binary requests
	DisableResponse      bool   // Disable response dumping
}

// gin Middleware for dump request and response
func Dumper(fn func(output *Output), options *Options) gin.HandlerFunc {
	if options == nil {
		options = &Options{}
	}
	return func(c *gin.Context) {

		buf := &bytes.Buffer{}
		buf.ReadFrom(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())) // Write body back
		reqBody := string(buf.Bytes())

		if options.TrimNewLineInRequest != "" {
			dct := http.DetectContentType(buf.Bytes())
			ct := strings.ToLower(c.Request.Header.Get("Content-Type"))

			if strings.HasPrefix(dct, "text") || ct == "application/json" || ct == "application/xml" {
				reqBody = strings.ReplaceAll(reqBody, "\n", options.TrimNewLineInRequest)
			}
		}

		var blw *bodyLogWriter
		if !options.DisableResponse {
			blw = &bodyLogWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}
			c.Writer = blw
		}

		c.Next() // Process request

		resp := ""
		if !options.DisableResponse {
			resp = blw.body.String()
		}

		fn(&Output{
			Request: &Request{
				Method:  c.Request.Method,
				Url:     c.Request.URL,
				Headers: c.Request.Header,
				Body:    reqBody,
			},
			Response: &Response{
				StatusCode: c.Writer.Status(),
				Body:       resp,
			},
		})
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
