package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/menuet"
	config "github.com/time-cop/timecop/pkg/config"
	database "github.com/time-cop/timecop/pkg/database"
)

func helloClock() {
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "Time " + time.Now().Format(":05"),
		})
		time.Sleep(time.Second)
	}
}

func humanDuration(duration float32) string {
	if duration < 1 {
		return fmt.Sprintf("%.0f secs", duration*60)
	}
	return fmt.Sprintf("%.0f mins", duration)
}

func menuItem(db *database.MemoryDatabase, task *database.Task) func() []menuet.MenuItem {
	return func() []menuet.MenuItem {
		items := []menuet.MenuItem{}

		title := task.Title

		items = append(items, menuet.MenuItem{
			Text: title,
		})

		items = append(items, menuet.MenuItem{
			Text: "snooze",
			Clicked: func() {
				db.SnoozeTask(task)
				db.Sort()
			},
		})

		items = append(items, menuet.MenuItem{
			Text: "finish",
			Clicked: func() {
				db.CompleteTask(task)
				db.Sort()
			},
		})

		return items
	}
}

func menuItems(db *database.MemoryDatabase) func() []menuet.MenuItem {
	return func() []menuet.MenuItem {
		items := []menuet.MenuItem{}
		db.Sort()
		currentTask := db.Tasks[db.CurrentTaskIndex]
		for _, task := range db.Incomplete() {
			title := task.Title
			if task == currentTask {
				title = fmt.Sprintf("... %s", title)
			}
			if task.SnoozeTimeLeft > 0 {
				timeLeft := humanDuration(task.SnoozeTimeLeft)
				title = fmt.Sprintf("%s ⏳%s", title, timeLeft)
			}
			item := menuet.MenuItem{
				Text: title,
				Clicked: func() {

				},
				Children: menuItem(db, task),
			}
			items = append(items, item)
		}

		return items
	}
}

func main() {
	go helloClock()
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
	menuet.App().Label = "com.github.pentaphobe.timecop"
	menuet.App().Children = menuItems(db)
	menuet.App().RunApplication()
}
