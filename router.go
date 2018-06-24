package dht

type Router struct {
	local   ID
	k       int
	buckets [BucketSize]Peers
}

func NewRouter(local ID, k int) (*Router, error) {
	r := new(Router)
	r.local = local
	r.k = k
	for i := range r.buckets {
		r.buckets[i] = Peers{}
	}
	return r, nil
}

func (r *Router) bucketIndex(id ID) int {
	var bi int
	if r.local == id {
		bi = len(r.buckets) - 1
	} else {
		bi = len(r.buckets) - r.local.Xor(id).BitLen()
	}
	return bi
}

func (r *Router) update(peer *Peer) {
	if r.local == peer.ID {
		return
	}
	bi := r.bucketIndex(peer.ID)
	bucket := r.buckets[bi]
	if len(bucket) == r.k {
		return
	}
	bucket[peer.ID] = peer
	return
}

func (r *Router) remove(id ID) {
	bi := r.bucketIndex(id)
	bucket := r.buckets[bi]
	delete(bucket, id)
}

func (r *Router) near(target ID, k int) []*Peer {
	result := make([]*Peer, 0, k)
	bi := r.bucketIndex(target)
	for i := bi; i < len(r.buckets); i++ {
		if len(result) >= k {
			break
		}
		bucket := r.buckets[i]
		result = append(result, bucket.List()...)
	}
	for i := bi - 1; i >= 0; i-- {
		if len(result) >= k {
			break
		}
		bucket := r.buckets[i]
		result = append(result, bucket.List()...)
	}
	result = NewPeers(result).Sort(target)
	if len(result) > k {
		result = result[:k]
	}
	return result
}

func (r *Router) length() int {
	length := 0
	for _, bucket := range r.buckets {
		length = length + len(bucket)
	}
	return length
}

func (r *Router) each(fn func(peer *Peer) bool) bool {
	for _, bucket := range r.buckets {
		for _, peer := range bucket {
			success := fn(peer)
			if success == false {
				return false
			}
		}
	}
	return true
}

func (r *Router) list() []*Peer {
	result := []*Peer{}
	for _, bucket := range r.buckets {
		result = append(result, bucket.List()...)
	}
	return result
}
