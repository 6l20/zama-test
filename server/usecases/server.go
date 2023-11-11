package usecases

import (
	"net/http"

	"github.com/6l20/zama-test/common/log"
	"github.com/6l20/zama-test/server"
	"github.com/6l20/zama-test/server/stores"
)

type ServerUseCases struct {
	server server.IServer
	logger log.Logger
	store stores.IStore
}

func NewServerUseCases(logger log.Logger, server server.IServer, store stores.IStore) *ServerUseCases {
return &ServerUseCases{
		server: server,
		logger: logger,
		store: store,
	}
}

func (s *ServerUseCases) HandleFileUpload() http.HandlerFunc {
	s.logger.Info("HandleFileUpload")
	return s.server.HandleFileUpload()
}

func (s *ServerUseCases) HandleFileRequest() http.HandlerFunc {
	s.logger.Info("HandleFileRequest")
	return s.server.HandleFileRequest()
}




