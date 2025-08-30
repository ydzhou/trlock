package trlock

import (
	"fmt"
	"os/exec"
	"strings"
)

func addInvalidPeersToPFRules(peers []*Peer) ([]byte, error) {
	cmdString := fmt.Sprintf("pfctl -t %s -T add", DefaultPFTable)
	for _, p := range peers {
		cmdString = fmt.Sprintf("%s %s", cmdString, p.IP)
	}
	args := strings.Split(cmdString, " ")
	cmd := exec.Command(args[0], args[1:]...)
	return runCmd(cmd)
}

func flushPFState() ([]byte, error) {
	cmd := exec.Command("pfctl", "-F", "state")
	return runCmd(cmd)
}

func resetPFRules() ([]byte, error) {
	cmd := exec.Command("pfctl", "-t", DefaultPFTable, "-T", "flush")
	return runCmd(cmd)
}

func runCmd(cmd *exec.Cmd) ([]byte, error) {
	var output []byte
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}
