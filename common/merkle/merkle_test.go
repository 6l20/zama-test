package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/6l20/zama-test/common/log/zap"
	"github.com/stretchr/testify/assert"
)

func TestMerkle1(t *testing.T) {

	zapLogger, _ := zap.NewLogger(&zap.Config{
		Level:  "debug",
		Format: "text",
	})
	manager := NewMerkleManager("/",zapLogger)

	dataBlocks := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	h0 := sha256.Sum256(dataBlocks[0])
	leafHash0 := hex.EncodeToString(h0[:])
	h1 := sha256.Sum256(dataBlocks[1])
	leafHash1 := hex.EncodeToString(h1[:])
	h2 := sha256.Sum256(dataBlocks[2])
	leafHash2 := hex.EncodeToString(h2[:])
	h3 := sha256.Sum256(dataBlocks[3])
	leafHash3 := hex.EncodeToString(h3[:])

	leafHashes := []string {
		leafHash0,
		leafHash1,
		leafHash2,
		leafHash3,
	}

	root := manager.PopulateTree(dataBlocks)

	assert.Equal(t, "58c89d709329eb37285837b042ab6ff72c7c8f74de0446b091b6a0131c102cfd", root.Hash)

	proof := manager.GenerateProof(1)

	verified := manager.VerifyProof(leafHashes[1], *proof, root.Hash)

	assert.Equal(t, true, verified)

	proof = manager.GenerateProof(0)

	verified = manager.VerifyProof(leafHashes[0], *proof, root.Hash)

	assert.Equal(t, true, verified)

}