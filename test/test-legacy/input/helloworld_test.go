package helloworld

import (
	"testing"
)

func TestSuccess(t *testing.T) {
	t.Log("will succeed")
}

func TestFailure(t *testing.T) {
	t.Log("will fail")
	t.Fail()
}
