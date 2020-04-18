package main

import (
	"fmt"

	"github.com/caseymrm/menuet"
	config "github.com/time-cop/timecop/pkg/config"
	database "github.com/time-cop/timecop/pkg/database"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	fmt.Printf("App paths: %#v\n", config.AppPaths)
	db := database.NewMemoryDatabase()
	// the following tasks should end up sorted
	db.AddTask(&database.Task{
		Title:      "Tidy room",
		IsComplete: true,
	})
	db.AddTask(&database.Task{
		Title:           "Sleep",
		LengthInMinutes: 1000,
		SnoozeTimeLeft:  20.0,
	})
	db.AddTask(&database.Task{
		Title:           "Procrastinate",
		LengthInMinutes: 1000,
		IsComplete:      true,
	})
	db.AddTask(&database.Task{
		Title:           "Work more on timecop",
		LengthInMinutes: 1000,
		Priority:        10,
	})
	snoozeTask := db.AddTask(&database.Task{
		Title:           "Work more on actual work",
		LengthInMinutes: 1000,
	})
	db.Sort()
	fmt.Printf("Tasks: %s\n", db.Tasks)
	fmt.Println("Snoozing a task")
	db.SnoozeTask(snoozeTask)
	db.Sort()
	fmt.Printf("Tasks: %s\n", db.Tasks)
	go helloClock(db)
	menuet.App().Label = "com.github.pentaphobe.timecop"
	menuet.App().Children = menuItems(db)
	menuet.App().RunApplication()
}
