package dht

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"sort"
)

type ID [IDSize]byte

func GenID() ID {
	id := ID{}
	_, err := rand.Reader.Read(id[:])
	if err != nil {
		panic(err)
	}
	return id
}

func FromID(raw []byte) ID {
	id := ID{}
	size := len(raw)
	if size != IDSize {
		return id
	}
	copy(id[:], raw)
	return id
}

func FromIDHex(hexstr string) ID {
	raw, err := hex.DecodeString(hexstr)
	if err != nil {
		return ID{}
	}
	return FromID(raw)
}

func (id ID) Xor(other ID) ID {
	var ret ID
	for i, b := range id {
		ret[i] = other[i] ^ b
	}
	return ret
}

func (id ID) BitLen() int {
	var num big.Int
	num.SetBytes(id[:])
	return num.BitLen()
}

func (id ID) Bytes() []byte {
	return id[:]
}

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

type IDList []ID

func (list IDList) Clone() IDList {
	result := make(IDList, len(list))
	copy(result, list)
	return result
}

func (list IDList) Sort(target ID) IDList {
	is := new(IDSort)
	is.Target = target
	is.List = list
	sort.Sort(is)
	return is.List
}

func (list IDList) ToSet() IDSet {
	result := IDSet{}
	for _, id := range list {
		result[id] = struct{}{}
	}
	return result
}

func (list IDList) Limit(size int) IDList {
	result := list.Clone()
	if len(result) > size {
		result = result[:size]
	}
	return result
}

type IDSet map[ID]struct{}

func (set IDSet) Clone() IDSet {
	result := IDSet{}
	for id, _ := range set {
		result[id] = struct{}{}
	}
	return result
}

func (set IDSet) Union(other IDSet) IDSet {
	result := set.Clone()
	for id, _ := range other {
		result[id] = struct{}{}
	}
	return result
}

func (set IDSet) Difference(other IDSet) IDSet {
	result := set.Clone()
	for id, _ := range other {
		_, exists := set[id]
		if exists == false {
			continue
		}
		delete(result, id)
	}
	return result
}

func (set IDSet) ToList() IDList {
	result := make(IDList, 0, len(set))
	for id, _ := range set {
		result = append(result, id)
	}
	return result
}

type IDSort struct {
	Target ID
	List   IDList
}

func (is *IDSort) Len() int {
	return len(is.List)
}

func (is *IDSort) Less(i, j int) bool {
	var a big.Int
	var b big.Int
	a.SetBytes(is.List[i].Xor(is.Target).Bytes())
	b.SetBytes(is.List[j].Xor(is.Target).Bytes())
	return a.Cmp(&b) == -1
}

func (is *IDSort) Swap(i, j int) {
	is.List[i], is.List[j] = is.List[j], is.List[i]
}
