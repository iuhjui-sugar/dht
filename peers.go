package dht

type Peer struct {
	ID   ID
	Addr string
}

type Peers map[ID]*Peer

func NewPeers(items ...*Peer) Peers {
	ps := Peers{}
	for _, item := range items {
		ps[item.ID] = item
	}
	return ps
}

func (ps Peers) Keys() IDList {
	result := IDList{}
	for id, _ := range ps {
		result = append(result, id)
	}
	return result
}

func (ps Peers) List() []*Peer {
	result := []*Peer{}
	for _, peer := range ps {
		result = append(result, peer)
	}
	return result
}

func (ps Peers) Sort(target ID) []*Peer {
	result := []*Peer{}
	keys := ps.Keys().Sort(target)
	for _, id := range keys {
		result = append(result, ps[id])
	}
	return result
}

func (ps Peers) Clone() Peers {
	result := Peers{}
	for id, peer := range ps {
		result[id] = peer
	}
	return result
}

func (ps Peers) Union(other Peers) Peers {
	result := ps.Clone()
	for id, peer := range other {
		result[id] = peer
	}
	return result
}

func (ps Peers) Extract(keys IDList) Peers {
	result := Peers{}
	for _, id := range keys {
		peer, exists := ps[id]
		if exists == false {
			continue
		}
		result[id] = peer
	}
	return result
}
