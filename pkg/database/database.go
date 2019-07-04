package database

import (
	"sort"
	"fmt"
	"time"
//	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"github.com/time-cop/timecop/pkg/config"
)

type Task struct {
	ID	string	`yaml:id`
	Title	string	`yaml:title`
	originalLengthInMinutes float32	`yaml:originalLength`
	LengthInMinutes	float32	`yaml:length`
	Priority	uint	`yaml:priority`
	SnoozeTimeLeft	float32	`yaml:snoozeTimeLeft`
	IsComplete	bool	`yaml:isComplete`
	CreatedAt	int64	`yaml:createdAt`
	UpdatedAt	int64	`yaml:updatedAt`
	CompletedAt	int64	`yaml:completedAt`
}

type TaskList []*Task

type DatabaseStore struct {
	Tasks	TaskList
	CurrentTaskIndex	uint
}

type Database interface {
	AddTask(task *Task) *Task
	SnoozeTask(task *Task)
	ChangeTaskLength(task *Task, lengthInMinutes float32)
	CompleteTask(task *Task)
	Sort()
}

type MemoryDatabase struct {
	DatabaseStore
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		DatabaseStore{
			Tasks: TaskList{
			},
			CurrentTaskIndex: 0,
		},
	}
}

func (db *MemoryDatabase) AddTask(task *Task) *Task {
	task.ID = uuid.New().String()
	if task.LengthInMinutes == 0 {
		task.LengthInMinutes = config.AppConfig.DefaultTaskLengthInMinutes
	}
	task.originalLengthInMinutes = task.LengthInMinutes
	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	db.Tasks = append(db.Tasks, task)
	return task
}

func (db *MemoryDatabase) ChangeTaskLength(task *Task, lengthInMinutes float32) {
	task.LengthInMinutes = lengthInMinutes
	task.UpdatedAt = time.Now().Unix()
}

func (db *MemoryDatabase) SnoozeTaskTime(task *Task, snoozeInMinutes float32) {
	task.SnoozeTimeLeft += snoozeInMinutes
	task.UpdatedAt = time.Now().Unix()
}

func (db *MemoryDatabase) SnoozeTask(task *Task) {
	db.SnoozeTaskTime(task, config.AppConfig.SnoozeInMinutes)
}

func (db *MemoryDatabase) CompleteTask(task *Task) {
	task.IsComplete = true
	task.CompletedAt = time.Now().Unix()
}

// CalculatePriority returns a number which _increases_ with task unimportance
func (task *Task) CalculatePriority() int64 {
	var priority int64 = 0
	if task.IsComplete {
		// TODO: arbitrary magic number is not scalable
		priority += 1000000
	}
	priority += int64(task.Priority)
	if task.SnoozeTimeLeft != 0.0 {
		// TODO: arbitrary magic number is not scalable
		priority += int64(task.SnoozeTimeLeft * 10000)
	}
	fmt.Printf("Task %s priority calculated as %d\n", task.ID, priority)
	return priority
}

func (taskList TaskList) String() string {
	out, err := yaml.Marshal(taskList)
	if err != nil {
		return fmt.Sprintf("can't cover tasklist to yaml: %#v", err)
	}
	return fmt.Sprintf("= Task List =\n%s\n==----==\n", string(out))
}

func (task Task) String() string {
	out, err := yaml.Marshal(task)
	if err != nil {
		return fmt.Sprintf("can't convert task to yaml: %#v", err)
	}
	return fmt.Sprintf("= Task =\n%s\n==----==\n", string(out))
}

func (db *MemoryDatabase) Sort() {
	sort.SliceStable(db.Tasks, func(i, j int) bool {
		return db.Tasks[i].CalculatePriority() < db.Tasks[j].CalculatePriority()
	})
}
