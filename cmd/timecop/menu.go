package main

import (
	"fmt"

	"github.com/caseymrm/menuet"
	database "github.com/time-cop/timecop/pkg/database"
)

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
		currentTask := db.CurrentTask()
		for _, task := range db.Incomplete() {
			thisTask := task
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
					db.SetCurrentTask(thisTask)
				},
				Children: menuItem(db, task),
			}
			items = append(items, item)
		}

		return items
	}
}
