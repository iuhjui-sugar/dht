package dht

import (
	"github.com/footstone-io/socker"
	"github.com/vmihailenco/msgpack"
)

type FindNodeReq struct {
	User   ID
	Target ID
}

type FindNodeResp struct {
	Items []*Peer
}

func (dht *DHT) findNodeHandler(msg *socker.MsgEvent) ([]byte, error) {
	req := new(FindNodeReq)
	err := msgpack.Unmarshal(msg.Body, &req)
	if err != nil {
		return nil, err
	}
	dht.update(&Peer{ID: req.User, Addr: msg.Raddr})
	// 查找本地距离最近的节点
	items := dht.router.near(req.Target, dht.router.k)
	return msgpack.Marshal(&FindNodeResp{
		Items: items,
	})
}

func (dht *DHT) findNode(addr string, target ID) ([]*Peer, error) {
	req := &FindNodeReq{
		User:   dht.ID(),
		Target: target,
	}
	reqb, err := msgpack.Marshal(req)
	if err != nil {
		return nil, err
	}
	respb, err := dht.socker.Query(addr, "findnode", reqb)
	if err != nil {
		return nil, err
	}
	resp := &FindNodeResp{}
	err = msgpack.Unmarshal(respb, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (dht *DHT) FindNode(target ID) ([]*Peer, error) {
	visit := IDSet{}
	peers := dht.router.near(target, K)
	keys := NewPeers(peers...).Keys()
	// 迭代查询
	for {
		unvisit := keys.ToSet().Difference(visit).ToList()
		if len(unvisit) == 0 {
			break
		}
		ch := make(chan Peers, len(unvisit))
		tq := NewTaskQueue(len(unvisit))
		tq.ExecGo(func(i int) {
			peer := peers[unvisit[i]]
			items, err := dht.findNode(peer.Addr, target)
			if err != nil {
				ch <- NewPeers()
			} else {
				ch <- NewPeers(items...)
			}
		})
		// 处理批量查询的结果
		for i := 0; i < len(unvisit); i++ {
			items := <-ch
			keys = keys.ToSet().Union(items.Keys().ToSet()).ToList()
			peers = peers.Union(items)
		}
		visit = visit.Union(unvisit.ToSet())
		keys = keys.Sort(target).Limit(K)
		peers = peers.Extract(keys)
	}
	return peers.Sort(target), nil
}
