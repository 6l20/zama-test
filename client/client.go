package client

import (
	"github.com/6l20/zama-test/common/log"
	"github.com/6l20/zama-test/common/merkle"

	"github.com/6l20/zama-test/client/config"
	"github.com/6l20/zama-test/client/stores"
)

type IClient interface {
	// UploadFile uploads a file to the server
	UploadFile(name string) error
	// DownloadFile downloads a file from the server
	DownloadFile() error
	// Generate Merkle Tree
	GenerateMerkleTree() error
	// Get Merkle Proof for a file
	GetMerkleProofForFile(fileNum int) (*merkle.Proof, error)
	// Verify Merkle Proof
	VerifyMerkleProof(leafHash string) bool
	// StoreMerkleRoot stores the merkle root locally
	StoreMerkleRoot() error
	// ReadMerkleRoot reads the merkle root from the local store
	ReadMerkleRoot() (string, error)
}

type Client struct {
	Logger log.Logger
	Config config.Config
	Store stores.ILocalStore
	merkleManager *merkle.MerkleManager
	RootFile string
	merkleTree *merkle.MerkleTree
	currentProof *merkle.Proof
}

func NewClient(logger log.Logger, config config.Config, store stores.ILocalStore) (*Client, error) {
	return &Client{
		Logger: logger,
		Config: config,
		Store: store,
		RootFile: config.RootFile,
		merkleManager: merkle.NewMerkleManager(logger.WithComponent("merkle")),
	}, nil
}

func (c *Client) UploadFile(name string) error {
	return UploadFile(name, c.Config.UploadURL)
}

func (c *Client) DownloadFile() error {
	return nil
}

func (c *Client) GetMerkleProof(leafIndex int) (*merkle.Proof, error) {
	return Proof(leafIndex, c.Config.ProofURL)
}

func (c *Client) GenerateMerkleTree() error {
	tree, err := c.merkleManager.BuildMerkleTreeFromFS("test/data")
	if err != nil {
		return err
	}
	c.merkleTree = tree
	return nil
}

func (c *Client) GetMerkleProofForFile(fileNum int) (*merkle.Proof, error) {
	return  Proof(fileNum, c.Config.ProofURL)
}

func (c *Client) VerifyMerkleProof(leafHash string) bool  {
	return c.merkleManager.VerifyProof(leafHash, *c.currentProof, c.merkleTree.Root.Hash)
}

func (c *Client) StoreMerkleRoot() error {
	return c.Store.StoreMerkleRoot("/tmp/" + c.RootFile, c.merkleTree.Root.Hash)
}

func (c *Client) ReadMerkleRoot() (string, error) {
	return c.Store.GetMerkleRoot("/tmp/" + c.RootFile)
}

