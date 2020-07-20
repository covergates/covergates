package web

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/web"
	"github.com/gin-gonic/gin"
)

// HandleIndex return HTML
func HandleIndex(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := bytes.NewBuffer([]byte{})
		html := web.MustLookup("/index.html")
		t, _ := template.New("index").Parse(string(html))
		t.Execute(buffer, config.Server.Base)
		data, _ := ioutil.ReadAll(buffer)
		c.Data(200, "text/html; charset=UTF-8", data)
	}
}
