package main

import (
	"github.com/thomasherstad/blog-aggregator/internal/config"
	"github.com/thomasherstad/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
