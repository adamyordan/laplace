package core

import "testing"

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
	testIP.IPAddress = "192.168.0.175"
	if err := testIP.CreateBarrierSession(); err != nil {
		t.Error()
	}
	if err := testIP.DeleteBarrierSession(); err != nil {
		t.Error()
	}
}
