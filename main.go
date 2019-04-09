package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chousemath/pomodoro-cli/dbjson"
	"github.com/chousemath/pomodoro-cli/noti"
	"github.com/chousemath/pomodoro-cli/pomodoro"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	resetHelp := "Indicates that you want the check count to be reset"
	resetChecks := flag.Bool("reset", false, resetHelp)
	resetChecksShort := flag.Bool("r", false, resetHelp)

	pomSessLenHelp := "Indicates how long you want this session to be"
	pomSessLen := flag.Int64("length", 0, pomSessLenHelp)
	pomSessLenShort := flag.Int64("l", 0, pomSessLenHelp)

	goalHelp := "The text content of your goal"
	goalText := flag.String("goal", "", goalHelp)
	goalTextShort := flag.String("g", "", goalHelp)

	isServerHelp := "Indicates whether or not you want to run the Pomodoro server"
	isServer := flag.Bool("server", false, isServerHelp)
	isServerShort := flag.Bool("s", false, isServerHelp)

	flag.Parse()

	db := dbjson.LoadDB()

	if *resetChecks || *resetChecksShort {
		db.Checks = 0
	}

	if *pomSessLen == 0 && *pomSessLenShort != 0 {
		*pomSessLen = *pomSessLenShort
	}
	if *pomSessLen == 0 {
		*pomSessLen = pomodoro.SessionLength
	}

	if *isServer || *isServerShort {
		r := mux.NewRouter()
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			setHeaderHTML(&w)

			var goals strings.Builder
			goals.WriteString("<h3>My Past Goals</h3>")
			goals.WriteString("<ul>")
			db.SortGoals()
			for _, goal := range db.GoalList {
				goals.WriteString("<li>")
				goals.WriteString("<b>")
				unixTimeUTC := time.Unix(goal.CompletedAt, 0) //gives unix time stamp in utc
				goals.WriteString(unixTimeUTC.Format(time.RFC3339))
				goals.WriteString("</b>")
				goals.WriteString(" - ")
				goals.WriteString(goal.Description)
				goals.WriteString("</li>")
			}
			goals.WriteString("</ul>")

			fmt.Fprintf(w, goals.String())
		}).Methods("GET")
		log.Fatal(http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r)))
	}

	// update the user on a 5 minute interval.
	for i := int64(5); i < *pomSessLen; i += 5 {
		go noti.SleepThenNotify(i, *pomSessLen)
	}

	db.NotifyAndSleep(*pomSessLen)

	if *goalText == "" && *goalTextShort != "" {
		*goalText = *goalTextShort
	}

	if err := db.CheckAndNotify(*goalText); err != nil {
		log.Fatal(fmt.Sprintf("Error checking and notifying: %s", err.Error()))
	}
	db.Save()
}

func setHeaderHTML(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
}
