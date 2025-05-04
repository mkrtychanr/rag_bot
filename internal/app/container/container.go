package container

import (
	tgapi "github.com/mkrtychanr/rag_bot/internal/gateway/tg_api"
	"github.com/mkrtychanr/rag_bot/internal/repository"
)

type Container struct {
	TgGateway  *tgapi.TgAPI
	Repository *repository.Repository
}
