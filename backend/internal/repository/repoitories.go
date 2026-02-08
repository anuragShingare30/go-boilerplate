package repository

import "github.com/anuragShingare30/go-boilerplate/internal/server"

// @dev This contains different db methods
type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
