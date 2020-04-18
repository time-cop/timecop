package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/menuet"
	database "github.com/time-cop/timecop/pkg/database"
)

func helloClock(db *database.MemoryDatabase) {
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: db.CurrentTask().Title + " " + time.Now().Format(":05"),
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
