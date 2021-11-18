// +build !windows

package drop_privilege

import (
	"testing"
)

func TestGetpwnam(t *testing.T) {
	pwd := Getpwnam("nobody")
	if pwd == nil {
		t.Error("getpwnam failed")
	}
}
