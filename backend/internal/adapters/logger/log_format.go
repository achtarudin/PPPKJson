package logger

import (
	"bytes"
	"time"

	"github.com/fatih/color"
)

type logFormat struct{}

func (writer *logFormat) Write(result []byte) (n int, err error) {
	color.NoColor = false

	var c *color.Color

	lower := bytes.ToLower(result)

	if bytes.Contains(lower, []byte("error")) || bytes.Contains(lower, []byte("failed")) {
		c = color.New(color.FgRed)
	} else if bytes.Contains(lower, []byte("info")) {
		c = color.New(color.FgBlue, color.Bold)
	} else if bytes.Contains(lower, []byte("warn")) {
		c = color.New(color.FgYellow, color.Bold)
	} else {
		c = color.New(color.FgGreen)
	}
	return c.Print(time.Now().UTC().Format("02/01/2006 15:04:05") + " " + string(result))
}
