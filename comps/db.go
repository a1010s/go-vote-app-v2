package comps

import (
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

var (
	DB       *badger.DB
	VotedIPs = make(map[string]bool)
	Mu       sync.Mutex
)

func InitDB() {
	var err error
	opts := badger.DefaultOptions("badger-db")
	DB, err = badger.Open(opts)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
}
