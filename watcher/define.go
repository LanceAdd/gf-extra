package watcher

import (
	"sync"
)

type SnapShot struct {
	sync.RWMutex
	CfgMap *map[string]any
}
