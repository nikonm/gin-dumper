package gin_dumper

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDumper(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		respBody             string
		expReqBody           string
		expRespBody          string
		trimNewLineInRequest string
	}{
		{
			name:                 "emptyOptions",
			requestBody:          "{\"test\": \"ok\",\n\"data\":123}",
			respBody:             `{"done": true}`,
			expReqBody:           "{\"test\": \"ok\",\n\"data\":123}",
			expRespBody:          `{"done": true}`,
			trimNewLineInRequest: "",
		},
		{
			name:                 "OptionsTrim",
			requestBody:          "{\"test\": \"ok\",\n\"data\":123}",
			respBody:             `{"done": true}`,
			expReqBody:           "{\"test\": \"ok\",   \"data\":123}",
			expRespBody:          `{"done": true}`,
			trimNewLineInRequest: "   ",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()

			var opt *Options
			if test.trimNewLineInRequest != "" {
				opt = &Options{TrimNewLineInRequest: test.trimNewLineInRequest}
			}

			router.Use(Dumper(func(output *Output) {
				require.Equal(t, test.expReqBody, output.Request.Body)
				require.Equal(t, test.expRespBody, output.Response.Body)
			}, opt))

			router.POST("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "%s", test.respBody)
			})

			ct := "application/json"
			req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(test.requestBody))
			req.Header.Set("Content-Type", ct)

			router.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
