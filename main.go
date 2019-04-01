package main

import (
	"time"

	"github.com/gen2brain/beeep"
)

func main() {
	if err := beeep.Notify(
		"Pomodoro timer started, work hard!",
		"Concentrate and get shit done Jo, you will be rewarded with a break.",
		"assets/clippy.png",
	); err != nil {
		panic(err)
	}

	// original pomodoro technique suggests a 25 min work cycle
	time.Sleep(25 * time.Minute)

	if err := beeep.Notify(
		"Time to take a walk!",
		"Make sure to make a check mark on your board.",
		"assets/clippy.png",
	); err != nil {
		panic(err)
	}
}
