package helpers

import (
	"log"
	"time"
)

// FUNCTION TO AID IN TIMING CODE BLOCKS

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
