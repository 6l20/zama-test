package stores

import (
	"os"
	"path/filepath"
)

type ILocalStore interface {
	// GetFile returns a file from the local store
	GetMerkleRoot(fileName string) (string, error)
	// StoreFile stores a file in the local store
	StoreMerkleRoot(filename, root string) error
}

// LocalStore implements ILocalStore interface
type LocalStore struct {
	baseDir  string // base directory to store the Merkle root files
	fileName string // name of the file to store the Merkle root
}

// NewLocalStore creates a new instance of LocalStore
func NewLocalStore(baseDir, fileName string) *LocalStore {
	return &LocalStore{
		baseDir:  baseDir,
		fileName: fileName,
	}
}

// GetMerkleRoot retrieves the Merkle root of a file from the local store
func (ls *LocalStore) GetMerkleRoot(fileName string) (string, error) {
	merkleFilePath := filepath.Join(ls.baseDir, ls.fileName+".merkle")
	data, err := os.ReadFile(merkleFilePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// StoreMerkleRoot stores the Merkle root of a file in the local store
func (ls *LocalStore) StoreMerkleRoot(filename, root string) error {
	merkleFilePath := filepath.Join(ls.baseDir, ls.fileName+".merkle")
	return os.WriteFile(merkleFilePath, []byte(root), 0644)
}
