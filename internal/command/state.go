package command

import (
	"github.com/KelvinJRosado/gator/internal/config"
	"github.com/KelvinJRosado/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}
