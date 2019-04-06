package noti

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
)

// Notify creates a desktop notification with a header and a body
func Notify(header, body string) {
	if err := beeep.Notify(header, body, "assets/clippy.png"); err != nil {
		panic(err)
	}
}

// SleepThenNotify sleeps for a certain amount of time, and then creates
// a desktop notification
func SleepThenNotify(sleepDuration, pomSessLen int64) {
	time.Sleep(time.Duration(sleepDuration) * time.Minute)
	Notify(
		"Keep it going!",
		fmt.Sprintf(
			"You have %d minutes left in this session.",
			pomSessLen-sleepDuration,
		),
	)
}
