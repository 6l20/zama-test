package usecases

import (
	"os"

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

func (c *ClientUseCases) UploadFilesFromDir(dirPath string) error {
	c.logger.Info("UploadFilesFromDir")
	files, err := os.ReadDir(dirPath)
    if err != nil {
        return err
    }

    for _, file := range files {
        if file.IsDir() {
			continue
		}
		err = c.Client.UploadFile(dirPath + "/" + file.Name())
		if err != nil {
			return err
		}
    }
	return nil
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
