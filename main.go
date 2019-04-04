package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
)

type goal struct {
	Description string `json:"Description"`
	CompletedAt int64  `json:"CompletedAt"`
}

type dbJSON struct {
	Checks         uint   `json:"Checks"`
	Sessions       uint   `json:"Sessions"`
	GoalComplete   uint   `json:"GoalComplete"`
	GoalIncomplete uint   `json:"GoalIncomplete"`
	GoalList       []goal `json:"GoalList"`
	FailureReasons []goal `json:"FailureReasons"`
	UpdatedAt      int64  `json:"UpdatedAt"`
}

const (
	// Yes indicates that a goal has been completed during the Pomodoro interval
	Yes string = "Yes, I finished my goal."
	// No indicates that I was unable to complete a goal during the interval
	No string = "No, I was unable to finish."
)

func main() {
	resetChecks := flag.Bool("reset", false, "Indicates that you want the check count to be reset")
	flag.Parse()

	db := loadDB()
	if *resetChecks {
		db.Checks = 0
	}

	go sleepThenNotify(5)
	db.notifyAndSleep()
	if err := db.checkAndNotify(); err != nil {
		log.Fatal(fmt.Sprintf("Error checking and notifying: %s", err.Error()))
	}
	db.save()
}

func notify(title, content string) {
	if err := beeep.Notify(title, content, "assets/clippy.png"); err != nil {
		panic(err)
	}
}

func loadDB() *dbJSON {
	dbFile, err := os.Open("./db.json")
	if err != nil {
		log.Fatal("Could not open the db file...")
	}
	defer dbFile.Close()
	dbBytes, err := ioutil.ReadAll(dbFile)
	if err != nil {
		log.Fatal("Could not convert db file to bytes...")
	}
	// initialize database configuration
	db := &dbJSON{}
	if err = json.Unmarshal(dbBytes, db); err != nil {
		log.Fatalf("Could not unmarshal db: %v", err)
	}
	return db
}

func (db *dbJSON) save() {
	db.UpdatedAt = time.Now().Unix()
	jsonData, err := json.Marshal(*db)
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Create("./db.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
}

func sleepThenNotify(sleepDuration int64) {
	time.Sleep(time.Duration(sleepDuration) * time.Minute)
	notify(
		"Keep it going!",
		fmt.Sprintf(
			"You have %d minutes left in this session.",
			25-sleepDuration,
		),
	)
}

func (db *dbJSON) notifyAndSleep() {
	notify(
		"Pomodoro timer started, work hard!",
		fmt.Sprintf(
			"Concentrate Jo, you currently have %d check%s.",
			db.Checks,
			pluralize(db.Checks),
		),
	)
	// original pomodoro technique suggests a 25 min work cycle
	time.Sleep(25 * time.Minute)
}

func (db *dbJSON) checkAndNotify() error {
	db.Checks++
	db.Sessions++
	breakTime := 5
	if db.Checks >= 4 {
		db.Checks = 0
		breakTime = 15
	}
	if err := db.checkGoal(); err != nil {
		return err
	}
	notify(
		"Time to take a walk!",
		fmt.Sprintf(
			"Take a %d minute break. You now have %d check%s.",
			breakTime,
			db.Checks,
			pluralize(db.Checks),
		),
	)
	return nil
}

func (db *dbJSON) checkGoal() error {
	answer, _, err := dlgs.List(
		"Goal Finished?",
		"Select an answer from the list:",
		[]string{Yes, No},
	)
	if err != nil {
		return err
	}

	switch answer {
	case Yes:
		db.GoalComplete++
		goalDesc, _, err := dlgs.Password("Description", "Describe your goal:")
		if err != nil {
			panic(err)
		}
		db.GoalList = append(
			db.GoalList,
			goal{
				Description: goalDesc,
				CompletedAt: time.Now().Unix(),
			},
		)
	case No:
		db.GoalIncomplete++
		failureDesc, _, err := dlgs.Password("Description", "Describe what went wrong:")
		if err != nil {
			panic(err)
		}
		db.FailureReasons = append(
			db.FailureReasons,
			goal{
				Description: failureDesc,
				CompletedAt: time.Now().Unix(),
			},
		)
	}
	return nil
}

func pluralize(checkCount uint) string {
	if checkCount == 1 {
		return ""
	}
	return "s"
}
