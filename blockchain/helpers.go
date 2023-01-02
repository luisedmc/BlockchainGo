package blockchain

import "log"

// HandleErrors is a helper function to handle errors.
func HandleErrors(err error) {
	if err != nil {
		log.Panic(err)
	}
}
