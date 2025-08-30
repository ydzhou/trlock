package trlock

import "time"

const DefaultConfigPath = "./trlock.conf"

const DefaultHostAddr = "127.0.0.1"
const DefaultHostPort = "9091"
const DefaultBlocklistPath = "./trblocklist"
const DefaultLogPath = "./trlock.log"
const DefaultPFTable = "trblocklist"

const DefaultInterval = 600 * time.Second
const DefaultResetInterval = 24 * time.Hour

const ConfigFileEnv = "TRLOCK_CONFIGFILE"
const LogFileEnv = "TRLOCK_LOGFILE"
const LogDebugEnv = "TRLOCK_LOGDEBUG"
