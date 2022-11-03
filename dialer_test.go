package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDroboInfo(t *testing.T) {
	assert := assert.New(t)

	d, err := getDroboInfo("169.254.6.109:5000")
	assert.NoError(err)
	assert.NotEmpty(d.Serial)

	t.Logf("%+v\n", d)
}
