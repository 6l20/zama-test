package usecases

import (
	"github.com/6l20/zama-test/client"
	"github.com/6l20/zama-test/common/log"
	"github.com/6l20/zama-test/common/merkle"
)

type ClientUseCases struct {
	Client client.IClient
	logger log.Logger
	
}

func NewClientUseCases(logger log.Logger, client client.IClient) *ClientUseCases {
	return &ClientUseCases{
		Client: client,
		logger: logger,
	}
}

func (c *ClientUseCases) UploadFile(name string) error {
	c.logger.Info("UploadFile")
	return c.Client.UploadFile(name)
}

func (c *ClientUseCases) DownloadFile() error {
	c.logger.Info("DownloadFile")
	return c.Client.DownloadFile()
}

func (c *ClientUseCases) GetMerkleProofForFile(fileNum int) (*merkle.Proof, error) {
	c.logger.Info("GetMerkleProofForFile")
	return c.Client.GetMerkleProofForFile(fileNum)
}

func (c *ClientUseCases) GenerateMerkleTree() error {
	c.logger.Info("GenerateMerkleTree")
	return c.Client.GenerateMerkleTree()
}

func (c *ClientUseCases) VerifyMerkleProof(leafHash string) bool {
	c.logger.Info("VerifyMerkleProof")
	return c.Client.VerifyMerkleProof(leafHash)
}

func (c *ClientUseCases) StoreMerkleRoot() error {
	c.logger.Info("StoreMerkleRoot")
	return c.Client.StoreMerkleRoot()
}

func (c *ClientUseCases) ReadMerkleRoot() (string, error) {
	c.logger.Info("ReadMerkleRoot")
	return c.Client.ReadMerkleRoot()
}
