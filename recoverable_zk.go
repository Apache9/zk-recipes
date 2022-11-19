package zk_recipes

import (
	"github.com/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)

type RecoverableZk struct {
	servers        []string
	sessionTimeout time.Duration
	stop           chan bool
	conn           *zk.Conn
	lock           *sync.RWMutex
}

func Connect(servers []string, sessionTimeout time.Duration) *RecoverableZk {
	rzk := &RecoverableZk{
		servers:        servers,
		sessionTimeout: sessionTimeout,
		stop:           make(chan bool, 1),
	}
	go func() {
		for {
			conn, watcher, err := zk.Connect(servers, sessionTimeout)
			if err != nil {
				time.Sleep(sessionTimeout)
				continue
			}
			rzk.lock.Lock()
			rzk.conn = conn
			rzk.lock.Unlock()
			for {
				select {
				case <-rzk.stop:
					log.Println("closed, quiting...")

					return
				case event := <-watcher:
					if event.State == zk.StateExpired {
						log.Println("session expired, try reconnecting...")
					}
				}
			}

		}
	}()

	return rzk
}
