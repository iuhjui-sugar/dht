package dht

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GenGroupID(size int) IDList {
	result := make(IDList, size)
	for i, _ := range result {
		result[i] = GenID()
	}
	return result
}

func TestGenID(t *testing.T) {
	id := GenID()
	assert.NotEqual(t, id, ID{})
}

func TestFromID(t *testing.T) {
	base := GenID()
	raw := base[:]
	result := FromID(raw)
	assert.Equal(t, result, base)

	raw = base[:2]
	result = FromID(raw)
	assert.Equal(t, result, ID{})
}

func TestIDListCopy(t *testing.T) {
	group1 := GenGroupID(3)
	group2 := group1.Clone()
	for i := 0; i < 3; i++ {
		assert.Equal(t, group1[i], group2[i])
	}
	group1[1] = GenID()
	assert.NotEqual(t, group1[1], group2[1])
}

func TestIDSetUnion(t *testing.T) {
	group1 := GenGroupID(3)
	group2 := GenGroupID(3)
	group1[0] = GenID()
	group2[0] = group1[0]
	set1 := group1.ToSet()
	set2 := group2.ToSet()
	set3 := set1.Union(set2)
	assert.Equal(t, len(set3), 5)
}
