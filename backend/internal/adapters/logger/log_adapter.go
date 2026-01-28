package logger

import (
	"log"
)

func New() {
	log.SetFlags(0)
	log.SetOutput(&logFormat{})
}
