// Package core
// We will implement the popular open source project barrier kvm
// To ensure we can connect clients keyboard and mouse to the server
// while ScreenShare is taking place
package core

import (
	"os/exec"
)

// Barrier It's preferred that the IP address used is a IPV6 address
type Barrier struct {
	NodeName    string
	IPAddress   string
	Mode        string
	Process     *exec.Cmd
}

// CreateBarrierSession Command to run "barrier.barrierc --debug INFO -f 192.168.0.175"
func (b *Barrier)CreateBarrierSession() error {

	// Checks if barrier client exists
	if err := DetectBarrier(); err != nil {
		return err
	}

	cmd := exec.Command("barrier.barrierc " ,b.IPAddress)
	err := cmd.Start()
	if err != nil {
		return err
	}
	// Saves the state of the command in the struct
	b.Process = cmd
	return nil
}

// DeleteBarrierSession Deletes barrier client session running
func (b *Barrier)DeleteBarrierSession() error {
	// Halts the process
	if err := b.Process.Wait(); err != nil {
		return err
	}
	return nil
}

// DetectBarrier This function ensures that the server has barrier client installed
func DetectBarrier() error {
	cmd := exec.Command("barrier.barrierc","-h")
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}