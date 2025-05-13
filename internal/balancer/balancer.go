package balancer

import (
	servFuncs "blcMod/internal/servers"
	"sync"
)

type Balancer struct {
	Current             int
	ReqCounter          []uint
	Servers             []*servFuncs.ServerStruct // Keep in mind, they are exist
	mu                  sync.Mutex                // for competitive req
	NumOfWorkingServers uint
}

func BalancerCreator(servers []*servFuncs.ServerStruct) *Balancer {
	return &Balancer{
		Current:             0,
		ReqCounter:          make([]uint, len(servers)),
		Servers:             servers,
		NumOfWorkingServers: 0,
	}
}

// Checking working servers at start of working.
// Later could be add gourutin with time duration and healthChecker
func (b *Balancer) CounterWorkingServers() {
	for _, val := range b.Servers {
		if val.ServerConfig.Readiness {
			b.NumOfWorkingServers++
		}
	}
}

// Checking all servers
func (b *Balancer) ChooseServer() *servFuncs.ServerStruct {
	// reqCounter occupied
	b.mu.Lock()
	defer b.mu.Unlock()

	// check if have working servers
	if len(b.Servers) == 0 || b.NumOfWorkingServers == 0 {
		return nil
	}

	minReq := 0
	for i := 1; i < int(b.NumOfWorkingServers); i++ {
		if b.ReqCounter[i] < b.ReqCounter[minReq] {
			minReq = i
		}
	}

	// increment counter of requests
	b.ReqCounter[minReq]++
	return b.Servers[minReq]
}
