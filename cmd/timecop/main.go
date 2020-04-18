package main

import (
	"fmt"
	"os"

	"github.com/caseymrm/menuet"
	"github.com/takama/daemon"
	config "github.com/time-cop/timecop/pkg/config"
	database "github.com/time-cop/timecop/pkg/database"
)

func main() {
	{
		srv, err := daemon.New(serviceName, description, dependencies...)
		if err != nil {
			errlog.Println("Error: ", err)
			os.Exit(1)
		}
		service := &Service{srv}
		status, err := service.Manage()
		if err != nil {
			errlog.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		fmt.Println(status)
		os.Exit(1)
	}

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
