package suite

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	println("hello world from test suite main")
	os.Exit(m.Run())
}

func TestSuccess(t *testing.T) {
	t.Log("integration test will succeed")
}

func TestFailure(t *testing.T) {
	t.Log("integration test will fail")
	t.Fail()
}
