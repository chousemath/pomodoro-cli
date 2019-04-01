package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
)

func main() {

	breakType := flag.String(
		"break",
		"short",
		"How long of a break will you take after this work cycle?",
	)

	flag.Parse()

	breakTime := 5
	if *breakType == "long" {
		breakTime = 15
	}

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
		fmt.Sprintf(
			"Take a %d minute break. Make sure to make a check mark on your board.",
			breakTime,
		),
		"assets/clippy.png",
	); err != nil {
		panic(err)
	}
}
