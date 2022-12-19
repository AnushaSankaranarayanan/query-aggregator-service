package probes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Liveness probe
func Liveness(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(c.Writer, "{ \"name\": \""+os.Getenv("NAME")+"\" }")
}
