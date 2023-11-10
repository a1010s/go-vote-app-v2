package comps

import (
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

func GetVotes(option string) int {
	var count int
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(option))
		if err == nil {
			err = item.Value(func(val []byte) error {
				count = int(val[0])
				return nil
			})
		}
		return err
	})
	if err != nil {
		count = 0
	}
	return count
}

func IndexHandler(c *gin.Context) {
	var option1Count, option2Count int
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(Option1))
		if err == nil {
			err = item.Value(func(val []byte) error {
				option1Count = int(val[0])
				return nil
			})
		}
		item, err = txn.Get([]byte(Option2))
		if err == nil {
			err = item.Value(func(val []byte) error {
				option2Count = int(val[0])
				return nil
			})
		}
		return err
	})
	if err != nil {
		option1Count, option2Count = 0, 0
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Question":     Question,
		"Option1":      Option1,
		"Option2":      Option2,
		"Option1Votes": option1Count,
		"Option2Votes": option2Count,
	})
}

func VoteHandler(c *gin.Context) {
	clientIP := c.ClientIP() // Get client's IP address
	Mu.Lock()
	defer Mu.Unlock()

	if VotedIPs[clientIP] {
		// User has already voted
		c.JSON(http.StatusForbidden, gin.H{"client IP": clientIP, "error": "You have already voted."})
		return
	}

	option := c.PostForm("vote")

	// Validate the selected option
	if option != Option1 && option != Option2 {
		log.Println("Invalid option:", option)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option"})
		return
	}

	// Update the votes
	key := []byte(option)
	err := DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == nil {
			var count int
			err = item.Value(func(val []byte) error {
				count = int(val[0])
				return nil
			})
			if err == nil {
				err = txn.Set(key, []byte{byte(count + 1)})
			}
		}
		return err
	})

	if err != nil {
		log.Println("Error updating votes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Mark the user's IP as voted
	VotedIPs[clientIP] = true

	// Retrieve the updated vote counts
	option1Count := GetVotes(Option1)
	option2Count := GetVotes(Option2)

	// Render the index.html template with the updated vote counts
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Question":     Question,
		"Option1":      Option1,
		"Option2":      Option2,
		"Option1Votes": option1Count,
		"Option2Votes": option2Count,
	})
}
