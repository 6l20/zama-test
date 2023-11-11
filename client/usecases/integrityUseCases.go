package usecases

type IIntegrtyUseCases interface {
	VerifyFile() error
	BuildMerkleRoot() error
}