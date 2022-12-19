package probes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbes(t *testing.T) {

	t.Run("Test probes : should pass", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		Liveness(c)

		if rr.Code != http.StatusOK {
			t.Errorf("Handler %s returned with error - got (%v) wanted (%v)", "liveness", rr.Code, http.StatusOK)
		}
	})

}
