package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Info("test hhh")
	Infof("test %v", 22)
	Error("test err")
	Errorf("error %v", 22)
}
