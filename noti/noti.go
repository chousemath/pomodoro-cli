package noti

import (
	"fmt"
	"sync"
	"time"

	"github.com/chousemath/pomodoro-cli/clear"
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
func SleepThenNotify(sleepDuration, pomSessLen int64, wg *sync.WaitGroup) {
	defer (*wg).Done()
	time.Sleep(time.Duration(sleepDuration) * time.Minute)
	msg := fmt.Sprintf(
		"You have %d minutes left in this session.",
		pomSessLen-sleepDuration,
	)
	Notify("Keep it going!", msg)
	clear.Execute()
	fmt.Println(msg)
}
