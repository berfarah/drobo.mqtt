package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadXML(t *testing.T) {
	assert := assert.New(t)

	b, err := os.ReadFile("./example/drobo-response.xml")
	assert.NoError(err)

	out, err := ReadXML(b)

	assert.Equal("drb141601a01445", out.Serial)
	assert.Equal("store", out.Name)
	assert.Equal(5917019996160, out.TotalCapacityProtected)
	assert.Equal([]string{"Normal"}, out.Statuses())
	assert.Equal("5.9 TB", out.TotalCapacity())
	assert.Equal("2.6 TB", out.UsedCapacity())
	assert.Equal("3.3 TB", out.FreeCapacity())

	assert.Len(out.Slots.Nodes, 6)

	slot1 := out.Slots.Nodes[0]
	assert.Equal("Hitachi  HUA72202", slot1.Make)
	assert.Equal("B9GMP7WF", slot1.Serial)
	assert.Equal(36, slot1.RotationalSpeed)
	assert.Equal("Green On", slot1.StatusString())
	assert.Equal("0 B", slot1.ManagedCapacity())
	assert.Equal("2.0 TB", slot1.PhysicalCapacity())
}
