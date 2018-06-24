package dht

import (
	"github.com/footstone-io/socker"
	"github.com/vmihailenco/msgpack"
)

type PingReq struct {
	User ID
}

type PingResp struct {
	ID ID
}

func (dht *DHT) pingHandler(msg *socker.MsgEvent) ([]byte, error) {
	req := new(PingReq)
	err := msgpack.Unmarshal(msg.Body, &req)
	if err != nil {
		return nil, err
	}
	dht.update(&Peer{ID: req.User, Addr: msg.Raddr})
	return msgpack.Marshal(&PingResp{
		ID: dht.ID(),
	})
}

func (dht *DHT) ping(addr string) (ID, error) {
	req := &PingReq{User: dht.ID()}
	reqb, err := msgpack.Marshal(req)
	if err != nil {
		return ID{}, err
	}
	respb, err := dht.socker.Query(addr, "ping", reqb)
	if err != nil {
		return ID{}, err
	}
	resp := &PingResp{}
	err = msgpack.Unmarshal(respb, &resp)
	if err != nil {
		return ID{}, err
	}
	return resp.ID, nil
}

func (dht *DHT) Ping(addr string) (*Peer, error) {
	id, err := dht.ping(addr)
	if err != nil {
		return nil, err
	}
	peer := &Peer{
		ID:   id,
		Addr: addr,
	}
	return peer, nil
}
