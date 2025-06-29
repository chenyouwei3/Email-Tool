package models

import (
	"net"
	"sync"
)

type Census struct {
	Mutex *sync.RWMutex
	Data  map[string]int
}

type BlackWhiteList struct {
	Mutex *sync.RWMutex
	Data  map[string]bool
}

type MxRecords struct {
	Mutex *sync.RWMutex
	Data  map[string][]*net.MX
}
