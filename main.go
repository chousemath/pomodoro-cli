package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chousemath/pomodoro-cli/dbjson"
	"github.com/chousemath/pomodoro-cli/noti"
	"github.com/chousemath/pomodoro-cli/pomodoro"
)

func main() {
	resetChecks := flag.Bool("reset", false, "Indicates that you want the check count to be reset")
	pomSessLen := flag.Int64("length", 0, "Indicates how long you want this session to be")
	goalText := flag.String("goal", "", "The text content of your goal")
	flag.Parse()

	db := dbjson.LoadDB()
	if *resetChecks {
		db.Checks = 0
	}
	if *pomSessLen == 0 {
		*pomSessLen = pomodoro.SessionLength
	}

	// update the user on a 5 minute interval
	for i := int64(5); i < *pomSessLen; i += 5 {
		go noti.SleepThenNotify(i, *pomSessLen)
	}

	db.NotifyAndSleep(*pomSessLen)
	if err := db.CheckAndNotify(*goalText); err != nil {
		log.Fatal(fmt.Sprintf("Error checking and notifying: %s", err.Error()))
	}
	db.Save()
}
