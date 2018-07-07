package meander_test

import (
	"github.com/avalchev94/go_blueprints/meander"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCostValues(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(int(meander.Cost1), 1)
	assert.Equal(int(meander.Cost2), 2)
	assert.Equal(int(meander.Cost3), 3)
	assert.Equal(int(meander.Cost4), 4)
	assert.Equal(int(meander.Cost5), 5)
}

func TestCostString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(meander.Cost1.String(), "$")
	assert.Equal(meander.Cost2.String(), "$$")
	assert.Equal(meander.Cost3.String(), "$$$")
	assert.Equal(meander.Cost4.String(), "$$$$")
	assert.Equal(meander.Cost5.String(), "$$$$$")
}

func TestParseCost(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(meander.Cost1, meander.ParseCost("$"))
	assert.Equal(meander.Cost2, meander.ParseCost("$$"))
	assert.Equal(meander.Cost3, meander.ParseCost("$$$"))
	assert.Equal(meander.Cost4, meander.ParseCost("$$$$"))
	assert.Equal(meander.Cost5, meander.ParseCost("$$$$$"))
}

func TestParseCostRange(t *testing.T) {
	assert := assert.New(t)

	l, err := meander.ParseCostRange("$$...$$$")
	assert.NoError(err)
	assert.Equal(l.From, meander.Cost2)
	assert.Equal(l.To, meander.Cost3)
	l, err = meander.ParseCostRange("$...$$$$$")
	assert.NoError(err)
	assert.Equal(l.From, meander.Cost1)
	assert.Equal(l.To, meander.Cost5)
}

func TestCostRangeString(t *testing.T) {
	assert := assert.New(t)
	r := meander.CostRange{
		From: meander.Cost2,
		To:   meander.Cost4,
	}
	assert.Equal("$$...$$$$", r.String())
}
