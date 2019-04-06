package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chousemath/pomodoro-cli/dbjson"
	"github.com/chousemath/pomodoro-cli/pomodoro"
)

func main() {
	resetChecks := flag.Bool("reset", false, "Indicates that you want the check count to be reset")
	flag.Parse()

	db := dbjson.LoadDB()
	if *resetChecks {
		db.Checks = 0
	}

	go sleepThenNotify(5)
	db.NotifyAndSleep(pomodoro.SessionLength)
	if err := db.CheckAndNotify(); err != nil {
		log.Fatal(fmt.Sprintf("Error checking and notifying: %s", err.Error()))
	}
	db.Save()
}
