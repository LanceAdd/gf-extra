package app

import (
	"sync"
)

type ApplicationConfig struct {
	sync.RWMutex
	AppName string
}

func GetAppName() string {
	appCfg.RLock()
	defer appCfg.RUnlock()
	return appCfg.AppName
}
