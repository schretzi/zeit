package z

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type Entry struct {
	ID      string    `json:"-"`
	Begin   time.Time `json:"begin,omitempty"`
	Finish  time.Time `json:"finish,omitempty"`
	Project string    `json:"project,omitempty"`
	Task    string    `json:"task,omitempty"`
	Notes   string    `json:"notes,omitempty"`
	User    string    `json:"user,omitempty"`

	SHA1 string `json:"-"`
}

func NewEntry(
	id string,
	begin string,
	finish string,
	project string,
	task string,
	user string) (Entry, error) {
	var err error

	newEntry := Entry{}

	newEntry.ID = id
	newEntry.Project = project
	newEntry.Task = task
	newEntry.User = user

	_, err = newEntry.SetBeginFromString(begin, time.Time{})
	if err != nil {
		return Entry{}, err
	}

	_, err = newEntry.SetFinishFromString(finish, time.Time{})
	if err != nil {
		return Entry{}, err
	}

	if id == "" && !newEntry.IsFinishedAfterBegan() {
		return Entry{}, errors.New("beginning time of tracking cannot be after finish time")
	}

	return newEntry, nil
}

func (entry *Entry) SetIDFromDatabaseKey(key string) error {
	splitKey := strings.Split(key, ":")

	if len(splitKey) < 3 || len(splitKey) > 3 {
		return errors.New("not a valid database key")
	}

	entry.ID = splitKey[2]
	return nil
}

func (entry *Entry) SetBeginFromString(begin string, contextTime time.Time) (time.Time, error) {
	var beginTime time.Time
	var err error

	if begin == "" {
		beginTime = time.Now()
	} else {
		beginTime, err = ParseTime(begin, contextTime)
		if err != nil {
			return beginTime, err
		}
	}

	entry.Begin = beginTime
	entry.secondsBegin()
	return entry.Begin, nil
}

func (entry *Entry) SetFinishFromString(finish string, contextTime time.Time) (time.Time, error) {
	var finishTime time.Time
	var err error

	if finish != "" {
		finishTime, err = ParseTime(finish, contextTime)
		if err != nil {
			return finishTime, err
		}
	}

	entry.Finish = finishTime
	entry.secondsFinish()
	return entry.Finish, nil
}

func (entry *Entry) IsFinishedAfterBegan() bool {
	return (entry.Finish.IsZero() || entry.Begin.Before(entry.Finish))
}

func (entry *Entry) GetOutputForTrack(isRunning bool, wasRunning bool) string {
	var outputPrefix string = ""
	var outputSuffix string = ""

	now := time.Now()
	trackDiffNow := now.Sub(entry.Begin)
	durationString := fmtDuration(trackDiffNow)

	if isRunning && !wasRunning {
		outputPrefix = "began tracking"
	} else if isRunning && wasRunning {
		outputPrefix = "tracking"
		outputSuffix = fmt.Sprintf(" for %sh", color.FgLightWhite.Render(durationString))
	} else if !isRunning && !wasRunning {
		outputPrefix = "tracked"
	}

	if entry.Task != "" && entry.Project != "" {
		return fmt.Sprintf("%s %s %s on %s%s\n", CharTrack, outputPrefix, color.FgLightWhite.Render(entry.Task), color.FgLightWhite.Render(entry.Project), outputSuffix)
	} else if entry.Task != "" && entry.Project == "" {
		return fmt.Sprintf("%s %s %s%s\n", CharTrack, outputPrefix, color.FgLightWhite.Render(entry.Task), outputSuffix)
	} else if entry.Task == "" && entry.Project != "" {
		return fmt.Sprintf("%s %s task on %s%s\n", CharTrack, outputPrefix, color.FgLightWhite.Render(entry.Project), outputSuffix)
	}

	return fmt.Sprintf("%s %s task%s\n", CharTrack, outputPrefix, outputSuffix)
}

func (entry *Entry) GetDuration() decimal.Decimal {
	duration := entry.Finish.Sub(entry.Begin)
	if duration < 0 {
		duration = time.Since(entry.Begin)
	}
	return decimal.NewFromFloat(duration.Hours())
}

func (entry *Entry) GetOutputForFinish() string {
	var outputSuffix string = ""

	trackDiff := entry.Finish.Sub(entry.Begin)
	taskDuration := fmtDuration(trackDiff)

	outputSuffix = fmt.Sprintf(" for %sh", color.FgLightWhite.Render(taskDuration))

	if entry.Task != "" && entry.Project != "" {
		return fmt.Sprintf("%s finished tracking %s on %s%s\n", CharFinish, color.FgLightWhite.Render(entry.Task), color.FgLightWhite.Render(entry.Project), outputSuffix)
	} else if entry.Task != "" && entry.Project == "" {
		return fmt.Sprintf("%s finished tracking %s%s\n", CharFinish, color.FgLightWhite.Render(entry.Task), outputSuffix)
	} else if entry.Task == "" && entry.Project != "" {
		return fmt.Sprintf("%s finished tracking task on %s%s\n", CharFinish, color.FgLightWhite.Render(entry.Project), outputSuffix)
	}

	return fmt.Sprintf("%s finished tracking task%s\n", CharFinish, outputSuffix)
}

func (entry *Entry) GetOutput(full bool) string {
	var output string = ""
	var entryFinish time.Time
	var isRunning string = ""

	if entry.Finish.IsZero() {
		entryFinish = time.Now()
		isRunning = "[running]"
	} else {
		entryFinish = entry.Finish
	}

	trackDiff := entryFinish.Sub(entry.Begin)
	taskDuration := fmtDuration(trackDiff)
	if !full {

		output = fmt.Sprintf("%s %s on %s from %s to %s (%sh) %s",
			color.FgGray.Render(entry.ID),
			color.FgLightWhite.Render(entry.Task),
			color.FgLightWhite.Render(entry.Project),
			color.FgLightWhite.Render(entry.Begin.Format(ExampleDateIso)),
			color.FgLightWhite.Render(entryFinish.Format(ExampleDateIso)),
			color.FgLightWhite.Render(taskDuration),
			color.FgLightYellow.Render(isRunning),
		)
	} else {
		output = fmt.Sprintf("%s\n   %s on %s\n   %sh from %s to %s %s\n\n   Notes:\n   %s\n",
			color.FgGray.Render(entry.ID),
			color.FgLightWhite.Render(entry.Task),
			color.FgLightWhite.Render(entry.Project),
			color.FgLightWhite.Render(taskDuration),
			color.FgLightWhite.Render(entry.Begin.Format(ExampleDateIso)),
			color.FgLightWhite.Render(entryFinish.Format(ExampleDateIso)),
			color.FgLightYellow.Render(isRunning),
			color.FgLightWhite.Render(strings.Replace(entry.Notes, "\n", "\n   ", -1)),
		)
	}

	return output
}

func (entry *Entry) secondsBegin() {
	if viper.GetBool("time.no-seconds") {
		entry.Begin = entry.Begin.Truncate(time.Duration(time.Minute))
	}
}

func (entry *Entry) secondsFinish() {
	if viper.GetBool("time.no-seconds") {
		entry.Finish = entry.Finish.Truncate(time.Duration(time.Minute))
	}
}

func GetFilteredEntries(entries []Entry, project string, task string, since time.Time, until time.Time) ([]Entry, error) {
	var filteredEntries []Entry

	for _, entry := range entries {
		if project != "" && GetIdFromName(entry.Project) != GetIdFromName(project) {
			continue
		}

		if task != "" && GetIdFromName(entry.Task) != GetIdFromName(task) {
			continue
		}

		if !since.IsZero() && !since.Before(entry.Begin) && !since.Equal(entry.Begin) {
			continue
		}

		if !until.IsZero() && !until.After(entry.Finish) && !until.Equal(entry.Finish) {
			continue
		}

		filteredEntries = append(filteredEntries, entry)
	}

	return filteredEntries, nil
}
