package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// @dev logic to connect the db
// @dev Database Pooling: Opening certain number of connections in-hand, application will be more efficient in performance


type Database struct {
	Pool *pgxpool.Pool // to store pool
	log *zerolog.Logger // to log db related info
}

// DatabaseTimeout is the timeout duration for database operations in seconds.
const DatabaseTimeout = 10