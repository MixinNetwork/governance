package externals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMixinRPC(t *testing.T) {
	assert := assert.New(t)

	nodes, err := ListAllNodes()
	assert.Nil(err)
	assert.LessOrEqual(485, len(nodes))

	tx, err := ReadTransaction("a1eb53c84b94f4cd2063cf7ed745d1f726123144dff03648868df44e9d317cfb")
	assert.Nil(err)
	assert.NotNil(tx)
	assert.Equal("81a154c41000000000000000000000000000000000", tx.Extra)
	assert.Equal("a1eb53c84b94f4cd2063cf7ed745d1f726123144dff03648868df44e9d317cfb", tx.Hash)
}
