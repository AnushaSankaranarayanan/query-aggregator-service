package swagger

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	openAPIFile            = "openapi.json"
	index                  = "index.html"
	contentLength          = "Content-Length"
	contentType            = "Content-Type"
	contentTypeHTML        = "text/html"
	contentTypeJavascript  = "application/javascript"
	contentTypeCSS         = "text/css"
	contentTypeOctetStream = "application/octet-stream"
)

var logger = logrus.StandardLogger()

//go:embed  assets

var f embed.FS

// Build - serves swagger files
func Build(c *gin.Context) {
	resource := resourceRequested(c)

	if resource == openAPIFile {
		http.ServeFile(c.Writer, c.Request, resource)
		return
	}

	//Get the resource requested
	file, _ := f.ReadFile("assets/" + resource)
	writeContentType(resource, c)

	if resource != index {
		serveOtherResources(file, c)
		return
	}

	serveIndex(file, c)
}

func resourceRequested(c *gin.Context) string {
	source := c.Request.URL.Path[1:]
	if source == "api/v1/openapi/" || source == "api/v1/openapi" {
		source = index
	}
	if strings.Contains(source, "/") {
		source = strings.SplitAfter(source, "/")[3]
	}
	return source
}

func writeContentType(source string, c *gin.Context) {
	switch filepath.Ext(source) {
	case ".html":
		c.Writer.Header().Set(contentType, contentTypeHTML)
	case ".js":
		c.Writer.Header().Set(contentType, contentTypeJavascript)
	case ".css":
		c.Writer.Header().Set(contentType, contentTypeCSS)
	default:
		c.Writer.Header().Set(contentType, contentTypeOctetStream)
	}
}

func serveIndex(staticFile []byte, c *gin.Context) {
	indexHTML := string(staticFile)
	logger.Trace(indexHTML)
	c.Writer.Header().Set(contentLength, strconv.Itoa(len(indexHTML)))
	fmt.Fprint(c.Writer, indexHTML)
}

func serveOtherResources(staticFile []byte, c *gin.Context) {
	c.Writer.Header().Set(contentLength, strconv.Itoa(len(staticFile)))
	c.Writer.Write(staticFile)
}
