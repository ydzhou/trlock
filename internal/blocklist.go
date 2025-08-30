package trlock

import (
	"fmt"
	"os"
)

func (t *Trlock) addInvalidPeersToBlocklist(peers []*Peer) {
	if len(peers) == 0 {
		return
	}

	file, err := os.OpenFile(t.config.BlocklistPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		t.log.Errorf("failed to open blocklist: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.log.Errorf("failed to close file: %v", err)
		}
	}()

	for _, p := range peers {
		if _, err = fmt.Fprintf(file, "%s:%s-%s\n", p.ClientName, p.IP, p.IP); err != nil {
			t.log.Errorf("failed to write blocklist: %v", err)
		}
	}

	t.log.Debugf("update blocklist: added %d IPs", len(peers))
}

func (t *Trlock) resetBlocklist() {
	file, err := os.OpenFile(t.config.BlocklistPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.log.Errorf("failed to open blocklist: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.log.Errorf("failed to close file: %v", err)
		}
	}()

	t.log.Debugf("reset blocklist")
}
