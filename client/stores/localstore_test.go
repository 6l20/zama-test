package stores

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetMerkleRoot tests the GetMerkleRoot method
func TestGetMerkleRoot(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "teststore")
	assert.NoError(t, err)

	defer os.RemoveAll(tmpDir) // clean up
	filename := "testfile"
	store := NewLocalStore(tmpDir, filename)
	
	merkleRoot := "testroot1234"

	// Pre-store a Merkle root for retrieval
	err = store.StoreMerkleRoot(filename, merkleRoot)
	require.NoError(t, err)

	// Test retrieving the Merkle root
	retrievedRoot, err := store.GetMerkleRoot(filename)
	assert.NoError(t, err)
	assert.Equal(t, merkleRoot, retrievedRoot)
}
