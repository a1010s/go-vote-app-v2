package main

import (
	"flag"
	"fmt"
	"log"

	components "github.com/a1010s/appmodules/comps"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.StringVar(&components.Option1, "option1", "Option 1", "Text for Option 1")
	flag.StringVar(&components.Option2, "option2", "Option 2", "Text for Option 2")
	flag.StringVar(&components.Question, "question", "Vote for Your Favorite", "Question for the voting poll")
	flag.Parse()

	components.InitDB() // Initialize the Badger database

	// Initialize keys with initial values if not present in the database
	err := components.DB.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(components.Option1))
		if err != nil && err == badger.ErrKeyNotFound {
			err = txn.Set([]byte(components.Option1), []byte{0})
		}

		_, err = txn.Get([]byte(components.Option2))
		if err != nil && err == badger.ErrKeyNotFound {
			err = txn.Set([]byte(components.Option2), []byte{0})
		}

		return err
	})

	if err != nil {
		log.Fatal("Error initializing keys:", err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", components.IndexHandler)
	r.POST("/vote", components.VoteHandler)

	fmt.Println("Server is running on :8099")
	r.Run(":8099")
}
