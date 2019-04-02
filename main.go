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
)

type dbJSON struct {
	Checks uint `json:"Checks"`
}

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

	db := loadDB()

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

	db.Checks++
	notify(
		"Time to take a walk!",
		fmt.Sprintf(
			"Take a %d minute break. You now have %d check%s.",
			breakTime,
			db.Checks,
			pluralize(db.Checks),
		),
	)

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

func pluralize(checkCount uint) string {
	if checkCount == 1 {
		return ""
	}
	return "s"
}
