package z

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Task struct {
	Name          string `json:"name,omitempty"`
	GitRepository string `json:"gitRepository,omitempty"`
}

func trackTask() {
	user := GetCurrentUser()

	runningEntryId, err := database.GetRunningEntryId(user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if runningEntryId != "" {
		log.Fatalf("%s a task is already running\n", CharTrack)
	}

	if project == "" && viper.GetString("project.default") != "" {
		project = viper.GetString("project.default")
	}

	if project == "" && viper.GetBool("project.mandatory") {
		log.Fatal("project is mandatory but missing")
	}

	if task == "" && viper.GetBool("task.mandatory") {
		log.Fatal("task is mandatory but missing")
	}

	newEntry, err := NewEntry("", begin, finish, project, task, user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if notes != "" {
		newEntry.Notes = notes
	}

	isRunning := newEntry.Finish.IsZero()

	_, err = database.AddEntry(user, newEntry, isRunning)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	fmt.Print(newEntry.GetOutputForTrack(isRunning, false))
}

func resumeTask(index int) {
	user := GetCurrentUser()

	entries, err := database.ListEntries(user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}
	lastEntry := entries[len(entries)-index]

	runningEntryId, err := database.GetRunningEntryId(user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if runningEntryId != "" {
		log.Fatalf("%s a task is already running\n", CharTrack)
	}

	project = lastEntry.Project
	task = lastEntry.Task

	newEntry, err := NewEntry("", begin, finish, project, task, user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if lastEntry.Notes != "" {
		newEntry.Notes = lastEntry.Notes
	}

	isRunning := newEntry.Finish.IsZero()

	_, err = database.AddEntry(user, newEntry, isRunning)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	fmt.Print(newEntry.GetOutputForTrack(isRunning, false))
}

func finishTask(mode int) {

	user := GetCurrentUser()

	runningEntryId, err := database.GetRunningEntryId(user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if runningEntryId == "" {
		log.Fatalf(ErrorString, CharError, err)
	}

	runningEntry, err := database.GetEntry(user, runningEntryId)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	tmpEntry, err := NewEntry(runningEntry.ID, begin, finish, project, task, user)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	if begin != "" {
		runningEntry.Begin = tmpEntry.Begin
	}

	if finish != "" {
		runningEntry.Finish = tmpEntry.Finish
	} else {
		runningEntry.Finish = time.Now()
	}

	if mode == FinishWithMetadata {
		finishTaskMetadata(user, &runningEntry, &tmpEntry)
	}

	if !runningEntry.IsFinishedAfterBegan() {
		fmt.Printf("%s %+v\n", CharError, "beginning time of tracking cannot be after finish time")
		os.Exit(1)
	}

	_, err = database.FinishEntry(user, runningEntry)
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	fmt.Print(runningEntry.GetOutputForFinish())
}

func finishTaskMetadata(user string, runningEntry *Entry, tmpEntry *Entry) {

	if project != "" {
		runningEntry.Project = tmpEntry.Project
	}

	if task != "" {
		runningEntry.Task = tmpEntry.Task
	}

	if notes != "" {
		runningEntry.Notes = fmt.Sprintf("%s\n%s", runningEntry.Notes, notes)
	}

	if runningEntry.Task != "" {
		task, err := database.GetTask(user, runningEntry.Task)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		taskGit(&task, runningEntry)
	}

}

func taskGit(task *Task, runningEntry *Entry) {
	if task.GitRepository != "" && task.GitRepository != "-" {
		stdout, stderr, err := GetGitLog(task.GitRepository, runningEntry.Begin, runningEntry.Finish)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		if stderr == "" {
			runningEntry.Notes = fmt.Sprintf("%s\n%s", runningEntry.Notes, stdout)
		} else {
			fmt.Printf("%s notes were not imported: %+v\n", CharError, stderr)
		}
	}
}
