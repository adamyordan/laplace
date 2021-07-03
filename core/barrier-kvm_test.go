package core

import (
	"testing"
	"time"
)

// To run this test ensure you have barrier
// installed
func TestDetectBarrier(t *testing.T) {
	err := DetectBarrier()
	if err != nil {
		t.Error()
	}
}

func TestBarrier_CreateBarrierSession_DeleteBarrierSession(t *testing.T) {
	var testIP Barrier
	// Change this with the test machine controlling the
	// keyboard and mouse
	testIP.IPAddress = "192.168.0.175"
	if err := testIP.CreateBarrierSession(); err != nil {
		t.Error()
	}
	// Create delay
	time.Sleep(5 * time.Second)

	if err := testIP.DeleteBarrierSession(); err != nil {
		t.Error()
	}
}
