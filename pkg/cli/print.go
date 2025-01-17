package cli

import (
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/colorize"
	"log"
)

const (
	errorPrefix   = "ERROR"
	warningPrefix = "WARNING"
	infoPrefix    = "INFO"
	successPrefix = "SUCCESS"
)

func MadeInBasliqLabs() {
	fmt.Printf(
		"made in %s",
		colorize.
			New("basliq labs\n").
			Modify(colorize.BrightBlue).
			Commit())
}

func Print(message string) {
	log.Printf("%s", message)
}

func Error(message string, err error) {
	log.Printf(
		"%s: %s\n\tDigest: %v",
		colorize.New(errorPrefix).Modify(colorize.Red).Commit(),
		message,
		colorize.New(err.Error()).Modify(colorize.BrightBlack).Modify(colorize.Underline).Commit(),
	)
}

func Warning(message string) {
	log.Printf(
		"%s: %s\n",
		colorize.New(warningPrefix).Modify(colorize.Yellow).Commit(),
		message,
	)
}

func Info(message string) {
	log.Printf(
		"%s: %s\n",
		colorize.New(infoPrefix).Modify(colorize.Cyan).Commit(),
		message,
	)
}

func Success(message string) {
	log.Printf(
		"%s: %s\n",
		colorize.New(successPrefix).Modify(colorize.Green).Modify(colorize.Bold).Commit(),
		message,
	)
}
