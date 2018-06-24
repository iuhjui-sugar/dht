package dht

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/footstone-io/socker"
)

type DHT struct {
	conf     *DHTConfig
	router   *Router
	socker   *socker.Socker
	shutdown chan chan struct{}
	mu       sync.Mutex
}

func NewDHT(conf *DHTConfig) (*DHT, error) {
	dht := new(DHT)
	dht.conf = conf
	dht.shutdown = make(chan chan struct{}, 1)
	err := conf.Verify()
	if err != nil {
		return nil, err
	}
	err = dht.makeRouter()
	if err != nil {
		return nil, err
	}
	err = dht.makeSocker()
	if err != nil {
		return nil, err
	}
	dht.bootstrap()
	go dht.keepalive()
	return dht, nil
}

func (dht *DHT) makeRouter() error {
	router, err := NewRouter(dht.conf.ID, K)
	if err != nil {
		return err
	}
	dht.router = router
	return nil
}

func (dht *DHT) makeSocker() error {
	laddr := fmt.Sprint(dht.conf.IP, ":", dht.conf.Port)
	sock, err := socker.NewSocker(laddr)
	if err != nil {
		return err
	}
	sock.Register("ping", dht.pingHandler)
	sock.Register("findnode", dht.findNodeHandler)
	dht.socker = sock
	return nil
}

func (dht *DHT) update(peer *Peer) {
	dht.mu.Lock()
	defer dht.mu.Unlock()
	dht.router.update(peer)
}

func (dht *DHT) remove(id ID) {
	dht.mu.Lock()
	defer dht.mu.Unlock()
	dht.router.remove(id)
}

func (dht *DHT) bootstrap() {
	tq := NewTaskQueue(len(dht.conf.Seeds))
	tq.ExecGo(func(i int) {
		peer, err := dht.Ping(dht.conf.Seeds[i])
		if err != nil {
			return
		}
		dht.update(peer)
	})
	items, err := dht.FindNode(dht.ID())
	if err != nil {
		return
	}
	for _, peer := range items {
		dht.update(peer)
	}
}

func (dht *DHT) keepalive() {
	ticker := time.NewTicker(time.Duration(dht.conf.KeepAlive) * time.Second)
	for {
		select {
		case <-ticker.C:
			peers := dht.router.list()
			tq := NewTaskQueue(len(peers))
			tq.ExecGo(func(i int) {
				peer := peers[i]
				_, err := dht.ping(peer.Addr)
				if err != nil {
					dht.remove(peer.ID)
				}
			})
		case done := <-dht.shutdown:
			done <- struct{}{}
			return
		}
	}
}

func (dht *DHT) ID() ID {
	return dht.conf.ID
}

func (dht *DHT) Addr() net.Addr {
	return dht.socker.Addr()
}

func (dht *DHT) Close() error {
	done := make(chan struct{}, 1)
	dht.shutdown <- done
	<-done
	err := dht.socker.Close()
	if err != nil {
		return err
	}
	return nil
}
