// dbjson contains all of the code that interacts directly with the json database

package dbjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/chousemath/pomodoro-cli/noti"
	"github.com/chousemath/pomodoro-cli/stredit"
	"github.com/gen2brain/dlgs"
)

// DBJSON represents a simple JSON database for this project
type DBJSON struct {
	Checks         uint   `json:"Checks"`
	Sessions       uint   `json:"Sessions"`
	GoalComplete   uint   `json:"GoalComplete"`
	GoalIncomplete uint   `json:"GoalIncomplete"`
	GoalList       []goal `json:"GoalList"`
	FailureReasons []goal `json:"FailureReasons"`
	UpdatedAt      int64  `json:"UpdatedAt"`
}

type goal struct {
	Description string `json:"Description"`
	CompletedAt int64  `json:"CompletedAt"`
}

const (
	// Yes indicates that a goal has been completed during the Pomodoro interval
	Yes string = "Yes, I finished my goal."
	// No indicates that I was unable to complete a goal during the interval
	No string = "No, I was unable to finish."
)

// LoadDB loads a json database into memory
func LoadDB() *DBJSON {
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
	db := &DBJSON{}
	if err = json.Unmarshal(dbBytes, db); err != nil {
		log.Fatalf("Could not unmarshal db: %v", err)
	}
	return db
}

// SortGoals sorts the goals in the database (most recent first)
func (db *DBJSON) SortGoals() {
	sort.Slice(db.GoalList, func(i, j int) bool {
		return db.GoalList[i].CompletedAt > db.GoalList[j].CompletedAt
	})
}

// Save records the user's progress by writing the user's state to a JSON file
func (db *DBJSON) Save() {
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

// CheckAndNotify creates a virtual check mark as per the Pomodoro
// technique, then notifies the user of their well-earned break
func (db *DBJSON) CheckAndNotify(goalText string) error {
	db.Checks++
	db.Sessions++
	breakTime := 5
	if db.Checks >= 4 {
		db.Checks = 0
		breakTime = 15
	}
	if err := db.checkGoal(goalText); err != nil {
		return err
	}
	noti.Notify(
		"Time to take a walk!",
		fmt.Sprintf(
			"Take a %d minute break. You now have %d check%s.",
			breakTime,
			db.Checks,
			stredit.Pluralize(db.Checks),
		),
	)
	return nil
}

func (db *DBJSON) checkGoal(goalText string) error {
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
		if goalText == "" {
			goalText, _, err = dlgs.Password("Description", "Describe your goal:")
			if err != nil {
				panic(err)
			}
		}
		db.GoalList = append(
			db.GoalList,
			goal{
				Description: goalText,
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
