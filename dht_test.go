package dht

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GenerateDHT(t *testing.T, seeds []string) *DHT {
	dc := GenerateDHTConfig()
	dc.Seeds = seeds
	dc.KeepAlive = 10
	dht, err := NewDHT(dc)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return dht
}

func TestDHTPing(t *testing.T) {
	dht1 := GenerateDHT(t, nil)
	dht2 := GenerateDHT(t, nil)
	node, err := dht2.Ping(dht1.Addr().String())
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(node.ID, node.Addr)
	assert.Equal(t, dht1.ID(), node.ID)
	assert.Equal(t, dht1.router.length(), 1)
	time.Sleep(time.Second)
}

/*
func TestDHTKeepalive(t *testing.T) {
	dht1 := GenerateDHT(t, nil)
	dht2 := GenerateDHT(t, []string{dht1.Addr().String()})
	time.Sleep(5 * time.Second)
	dht2.Close()
	time.Sleep(30 * time.Second)
	assert.Equal(t, dht1.router.length(), 0)
	time.Sleep(time.Second)
}
*/

func TestDHTSmallFindNode(t *testing.T) {
	dht1 := GenerateDHT(t, nil)
	dht2 := GenerateDHT(t, nil)
	dht3 := GenerateDHT(t, nil)
	_, err := dht2.Ping(dht1.Addr().String())
	assert.NoError(t, err)
	nodes, err := dht3.findNode(dht1.Addr().String(), dht2.ID())
	assert.NoError(t, err)
	assert.Equal(t, nodes[0].ID, dht2.ID())
}

func TestDHTBigFindNode(t *testing.T) {
	dht1 := GenerateDHT(t, nil)
	//dht2 := GenerateDHT(t, []string{dht1.Addr().String()})
	dht3 := GenerateDHT(t, []string{dht1.Addr().String()})
	dht4 := GenerateDHT(t, []string{dht1.Addr().String()})
	time.Sleep(time.Second)
	nodes, err := dht4.FindNode(dht3.ID())
	assert.NoError(t, err)
	assert.Equal(t, nodes[0].ID, dht3.ID())
}

func TestDHTBatchBigFindNode(t *testing.T) {
	dhts := []*DHT{}
	for i := 0; i < 100; i++ {
		seeds := []string{}
		if len(dhts) > 0 {
			j := rand.Intn(len(dhts))
			seeds = append(seeds, dhts[j].Addr().String())
		}
		dht := GenerateDHT(t, seeds)
		dhts = append(dhts, dht)
	}
	time.Sleep(time.Second)

	sample := []int{}
	for len(sample) < 3 {
		k := rand.Intn(len(dhts))
		sample = append(sample, k)
	}

	sample1 := dhts[sample[0]]
	//sample2 := dhts[sample[1]]
	sample3 := dhts[sample[2]]

	nodes, err := sample1.FindNode(sample3.ID())
	assert.NoError(t, err)
	assert.Equal(t, sample3.ID(), nodes[0].ID)
}
