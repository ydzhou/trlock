package trlock

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hekmon/transmissionrpc/v3"
	logrus "github.com/sirupsen/logrus"
)

type Trlock struct {
	config         *Config
	client         *transmissionrpc.Client
	log            *logrus.Logger
	invalidPeerMap map[string]*Peer
}

func (t *Trlock) Setup() {
	configFilePath, logFilePath, logDebug := loadEnv()

	t.log = logrus.New()
	if logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		t.log.Out = logFile
	} else {
		t.log.Warnf("failed to setup log file: %v", err)
	}
	if logDebug {
		t.log.Level = logrus.DebugLevel
	}

	t.loadConfig(configFilePath)

	t.invalidPeerMap = map[string]*Peer{}

	if client, err := setupClient(t.config.HostAddr, t.config.HostPort); err != nil {
		t.log.Fatalf("failed to connect to transmission host: %v", err)
	} else {
		t.client = client
	}

	t.log.Info("loaded config")

	t.Reset()
}

func (t *Trlock) Run() {

	ctx := context.Background()

	interval, err := time.ParseDuration(t.config.Interval)
	if err != nil {
		interval = DefaultInterval
	}
	resetInterval, err := time.ParseDuration(t.config.ResetInterval)
	if err != nil {
		resetInterval = DefaultResetInterval
	}
	ticker := time.NewTicker(interval)
	resetTicker := time.NewTicker(resetInterval)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	t.log.Infof("running on host: %s:%s at interval: %s", t.config.HostAddr, t.config.HostPort, interval)

	for {
		select {
		case <-ticker.C:
			t.Validate(ctx)

		case <-resetTicker.C:
			t.Reset()

		case <-exit:
			return
		}
	}
}

func (t *Trlock) Validate(ctx context.Context) {
	invalidPeers := t.getInvalidPeers(ctx)
	newPeers := t.getInvalidPeersToAction(invalidPeers)

	if len(newPeers) == 0 {
		t.log.Info("no new invalid peer found")
		return
	}

	if t.config.BlocklistEnabled {
		t.addInvalidPeersToBlocklist(newPeers)
		_, err := t.client.BlocklistUpdate(ctx)
		if err != nil {
			t.log.Errorf("failed to update blocklist: %v", err)
		}
		t.log.Debugln("update transmission blocklist")
	}

	if t.config.PfEnabled {
		out, err := addInvalidPeersToPFRules(newPeers)
		if err != nil {
			t.log.Errorf("failed to update pf rule: %v", err)
		}
		t.log.Debugf("update pf rule: %s", string(out))
		out, err = flushPFState()
		if err != nil {
			t.log.Errorf("failed to flush pf state: %v", err)
		}
		t.log.Debugf("flush pf state: %s", string(out))
	}
	t.log.Infof("validated %d invalid peers and applied actions", len(newPeers))
}

func (t *Trlock) Reset() {
	if t.config.BlocklistEnabled {
		t.resetBlocklist()
	}
	if t.config.PfEnabled {
		out, err := resetPFRules()
		if err != nil {
			t.log.Errorf("failed to reset pf rule: %v", err)
		} else {
			t.log.Debugf("reset pf rule: %s", string(out))
		}
	}

	t.invalidPeerMap = map[string]*Peer{}
	t.log.Info("reset invalid peers")
}

func (t *Trlock) getInvalidPeers(ctx context.Context) []*Peer {
	var invalidPeers []*Peer
	torrents, err := getAllTorrents(t.client, ctx)
	if err != nil {
		t.log.Errorf("failed to get torrents: %v", err)
		return nil
	}

	peers := getAllPeers(torrents)

	for _, peer := range peers {
		if isPeerValid(peer, t.config) {
			continue
		}
		invalidPeers = append(invalidPeers, peer)
		t.log.Debugf("invalid peer: %s %s\n", peer.IP, peer.ClientName)
	}

	return invalidPeers
}

func (t *Trlock) getInvalidPeersToAction(peers []*Peer) []*Peer {
	t.log.Debugf("%d invalid peers already taken action", len(t.invalidPeerMap))

	newInvalidPeers := []*Peer{}
	for _, p := range peers {
		if _, exist := t.invalidPeerMap[p.IP]; !exist {
			newInvalidPeers = append(newInvalidPeers, p)
			t.invalidPeerMap[p.IP] = p
		}
	}

	return newInvalidPeers
}
