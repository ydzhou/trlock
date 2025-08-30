package trlock

import (
	"encoding/json"
	"os"
	"runtime"
)

type Config struct {
	HostAddr           string     `json:"host,omitempty"`
	HostPort           string     `json:"port,omitempty"`
	Interval           string     `json:"interval,omitempty"`
	ResetInterval      string     `json:"reset_interval,omitempty"`
	StrictAllowEnabled bool       `json:"strict_allow_enabled,omitempty"`
	PfEnabled          bool       `json:"pf_enabled,omitempty"`
	BlocklistEnabled   bool       `json:"blocklist_enabled,omitempty"`
	BlocklistPath      string     `json:"blocklist_path,omitempty"`
	Allowlist          ListConfig `json:"allowlist"`
	Blocklist          ListConfig `json:"blocklist"`
}

type ListConfig struct {
	Client []string `json:"client,omitempty"`
}

func (t *Trlock) loadConfig(path string) {
	config := Config{
		HostAddr:           DefaultHostAddr,
		HostPort:           DefaultHostPort,
		Interval:           DefaultInterval.String(),
		ResetInterval:      DefaultResetInterval.String(),
		StrictAllowEnabled: false,
		PfEnabled:          false,
		BlocklistEnabled:   true,
		BlocklistPath:      DefaultBlocklistPath,
	}

	file, err := os.Open(path)
	if err != nil {
		t.log.Fatalf("failed to open config file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.log.Errorf("failed to close file: %v", err)
		}
	}()
	t.log.Debugf("open config file: %s", path)

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&config); err != nil {
		t.log.Fatalf("failed to load config: %v", err)
	}

	t.config = &config

	if t.config.PfEnabled && !isFreebsd() {
		t.log.Warnf("pf not supported in %s, set pf_enabled to false", runtime.GOOS)
	}
	t.config.PfEnabled = t.config.PfEnabled && isFreebsd()
}

func loadEnv() (string, string, bool) {
	configFilePath := os.Getenv(ConfigFileEnv)
	if configFilePath == "" {
		configFilePath = DefaultConfigPath
	}
	logFilePath := os.Getenv(LogFileEnv)
	if logFilePath == "" {
		logFilePath = DefaultLogPath
	}
	logDebug := os.Getenv(LogDebugEnv) != "0"
	return configFilePath, logFilePath, logDebug
}
