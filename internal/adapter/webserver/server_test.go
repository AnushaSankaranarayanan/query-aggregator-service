package webserver

import (
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("TestServer should fail(invalid port)", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "999999999")
		server := NewServer(Services{})
		err := server.Run()
		if err != nil && err.Error() != "listen tcp: address 999999999: invalid port" {
			t.Errorf("TestServer should fail(invalid port) expected(%v) got (%v)", "listen tcp: address 999999999: invalid port", err)
		}

	})
}
