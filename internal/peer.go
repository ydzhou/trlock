package trlock

import (
	"strings"

	"github.com/hekmon/transmissionrpc/v3"
)

type Peer struct {
	IP         string
	ClientName string
}

func getAllPeers(torrents []transmissionrpc.Torrent) []*Peer {
	peers := []*Peer{}
	peerIPMap := map[string]bool{}

	for _, torrent := range torrents {
		for _, p := range torrent.Peers {
			_, exist := peerIPMap[p.Address]
			if !exist {
				peers = append(peers, &Peer{IP: p.Address, ClientName: p.ClientName})
				peerIPMap[p.Address] = true
			}
		}
	}

	return peers
}

func isPeerValid(peer *Peer, config *Config) bool {
	return isPeerClientValid(peer, config)
}

func isPeerClientValid(peer *Peer, config *Config) bool {
	clientName := strings.TrimSpace(strings.ToLower(peer.ClientName))
	if len(clientName) == 0 {
		return false
	}
	if !config.StrictAllowEnabled {
		for _, c := range config.Blocklist.Client {
			if strings.Contains(clientName, c) {
				return false
			}
		}
		return true
	} else {
		for _, c := range config.Allowlist.Client {
			if strings.Contains(clientName, c) {
				return true
			}
		}
	}
	return false
}
