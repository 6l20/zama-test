package minio

import (
	"github.com/6l20/zama-test/common/log"
	"github.com/minio/minio-go"
)

const(
	endpoint = "minio:9000"
    accessKeyID = "user"
    secretAccessKey = "SecurePasswd"
)


// Minio is a struct that implements the IMinio interface.
type Minio struct {
	// Logger is the logger used by the Minio struct.
	Logger log.Logger
	// Config is the config used by the Minio struct.
	Config *Config
	// Client is the client used by the Minio struct.
	Client *minio.Client
	// Bucket is the bucket used by the Minio struct.
	Bucket string
}

// NewMinio returns a new Minio struct.
func NewMinio(logger log.Logger, config *Config) (*Minio, error) {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		return nil, err
	}

	return &Minio{
		Logger: logger,
		Config: config,
		Client: client,
		Bucket: config.Bucket,
	}, nil
}

// UploadFile uploads a file to the Minio server.
func (m *Minio) UploadFile(name string) error {
	m.Client.FPutObject(m.Bucket, name, name, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return nil
}

// DownloadFile downloads a file from the Minio server.
func (m *Minio) DownloadFile() error {
	return nil
}

// UploadMerkleTree uploads a merkle tree to the Minio server.
func (m *Minio) UploadMerkleTree() error {
	return nil
}

// DownloadMerkleTree downloads a merkle tree from the Minio server.
func (m *Minio) DownloadMerkleTree() error {
	return nil
}

// Connect connects to the Minio server.
func (m *Minio) Connect() error {
	err := m.Client.MakeBucket( m.Bucket, "")
        if err != nil {
                // Check to see if we already own this bucket (which happens if you run this twice)
                exists, errBucketExists := m.Client.BucketExists( m.Bucket)
                if errBucketExists == nil && exists {
                        m.Logger.Debug("We already own %s\n", m.Bucket)
                } else {
					m.Logger.Fatal(err.Error())
                }
        } else {
                m.Logger.Info("Successfully created %s\n", m.Bucket)
        }
		return nil
}

// Disconnect disconnects from the Minio server.
func (m *Minio) Disconnect() error {
	return nil
}