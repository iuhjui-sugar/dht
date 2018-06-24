package dht

import (
	"errors"
	"net"
)

type DHTConfig struct {
	ID        ID
	IP        string
	Port      uint16
	KeepAlive uint32
	Seeds     []string
}

func GenerateDHTConfig() *DHTConfig {
	dc := new(DHTConfig)
	dc.ID = GenID()
	dc.IP = "0.0.0.0"
	dc.Port = 0
	dc.KeepAlive = 5 * 60
	return dc
}

func (dc *DHTConfig) Verify() error {
	if dc == nil {
		return errors.New("DHTConfig Not Init")
	}
	if net.ParseIP(dc.IP) == nil {
		return errors.New("DHTConfig IP Format Error")
	}
	if dc.KeepAlive < 10 {
		return errors.New("DHTConfig minimum time for keep alive is 10s")
	}
	return nil
}
