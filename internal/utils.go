package trlock

import "runtime"

func isFreebsd() bool {
	return runtime.GOOS == "freebsd"
}
