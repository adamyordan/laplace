// Package core
// We will implement the popular open source project barrier kvm
// To ensure we can connect clients keyboard and mouse to the server
// while ScreenShare is taking place
package core

import (
	"laplace/config"
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
	//Checks if barrier client exists

	if err := DetectBarrier(); err != nil {
		return err
	}

	//Get username from config file
	configResp, err := config.ConfigInit()
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo","-u",configResp.SystemUsername,"barrier.barrierc","-f","--debug", "DEBUG" ,"--log", "/tmp/barrier.log",b.IPAddress)

	// USE THE FOLLOWING TO DEBUG
	//cmdReader, err := cmd.StdoutPipe()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	//	return
	//}
	//
	//// the following is used to print output of the command
	//// as it makes progress...
	//scanner := bufio.NewScanner(cmdReader)
	//go func() {
	//	for scanner.Scan() {
	//		fmt.Printf("%s\n", scanner.Text())
	//		//
	//		// TODO:
	//		// send output to server
	//	}
	//}()
	//
	//if err := cmd.Start(); err != nil {
	//	return err
	//}

	if err := cmd.Start(); err != nil {
		return err
	}

	println(cmd.Path)


	// Saves the state of the command in the struct
	b.Process = cmd
	return nil
}

// DeleteBarrierSession Deletes barrier client session running
func (b *Barrier)DeleteBarrierSession() error {
	// Halts the process
	cmd := exec.Command("pkill" ,"barrierc")
	if err := cmd.Run(); err != nil {
	      return err
	}

	return nil
}

// DetectBarrier This function ensures that the server has barrier client installed
func DetectBarrier() error {
	_, err := exec.LookPath("barrier.barrierc")
	if err != nil {
		return err
	}
	return nil
}