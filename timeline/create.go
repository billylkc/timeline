package timeline

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/now"
)

//go:embed templates/*
var templates embed.FS

// ParseFile helps to parse differnt file format to the Task structs
func ParseFile(r io.Reader, filetype string) ([]Task, error) {
	var tasks []Task

	switch filetype {
	case "json":
		return parseJSON(r)

	case "csv":
		return parseCSV(r)

	default:
		return tasks, fmt.Errorf("Format is not supported - %s \n", filetype)
	}

	return tasks, nil
}

// CreateTimeline creates timeline file with the input settings
func CreateTimeline(w io.Writer, tasks []Task) error {

	t, err := template.ParseFS(templates, "templates/timeline.gohtml")
	if err != nil {
		return err
	}

	// Put summary and timeline into the same struct before parsing
	var (
		output  Output
		done    []Task
		ongoing []Task
	)

	// compare date for finished/ongoing tasks
	today := time.Now()
	for _, t := range tasks {
		if t.EndD.Before(today) {
			done = append(done, t)
		} else {
			ongoing = append(ongoing, t)
		}
	}

	output.Done = done
	output.Ongoing = ongoing
	output.Timeline = tasks

	err = t.Execute(w, output)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(w)
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func parseJSON(r io.Reader) ([]Task, error) {
	var all AllTasks
	var tasks []Task

	// Unmarshal JSON
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return tasks, err
	}
	err = json.Unmarshal(b, &all)
	if err != nil {
		return tasks, err
	}

	// From the Tasks node
	tasks = all.Tasks
	for i, _ := range tasks { // Convert string to javascript object
		start := convertDate(string(tasks[i].Start))
		end := convertDate(string(tasks[i].End))

		tasks[i].Start, _ = convertJSDate(start)
		tasks[i].End, _ = convertJSDate(end)

		tasks[i].StartD = start
		tasks[i].EndD = end
		tasks[i].Duration = calcDuration(start, end)
	}

	return tasks, nil
}

func parseCSV(r io.Reader) ([]Task, error) {
	var tasks []Task

	// Unmarshal to io.Reader
	err := gocsv.Unmarshal(r, &tasks)
	if err != nil {
		return tasks, err
	}

	for i, _ := range tasks { // Convert string to javascript object
		start := convertDate(string(tasks[i].Start))
		end := convertDate(string(tasks[i].End))

		tasks[i].Start, _ = convertJSDate(start)
		tasks[i].End, _ = convertJSDate(end)

		tasks[i].StartD = start
		tasks[i].EndD = end
		tasks[i].Duration = calcDuration(start, end)
	}

	return tasks, nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// convertJSDate converts a date object to javascript object
// And stupid javascript added one month somehow
func convertJSDate(t time.Time) (template.JS, error) {
	var (
		temp string
		res  template.JS
	)

	// handle javascript added month
	t = t.AddDate(0, -1, 0) // Subtract one month
	d := t.Format("2006-01-02")

	splits := strings.Split(string(d), "-")
	if len(splits) != 3 {
		return "", fmt.Errorf("Invalid date format. Should be in YYYY-mm-dd format.\nInput - %s \n.", d)
	}
	temp = fmt.Sprintf("new Date(%s, %s, %s)", splits[0], splits[1], splits[2])
	res = template.JS(temp)
	return res, nil
}

// convertDate convert string in YYYY-MM-DD format to time.Time
func convertDate(d string) time.Time {
	t, _ := now.Parse(d)
	return t
}

// calcDuration calculates no of working days (5 days a week, no handling on holidays) between start and end
func calcDuration(start, end time.Time) int {

	if end.Before(start) || end.Equal(start) {
		return 0
	}

	days := 1
	t := start
	for {
		if days >= 365 {
			return days // early stop condition. Limit to 1 year
		}

		if t.Equal(end) {
			return days
		}
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			days++
		}
		t = t.AddDate(0, 0, 1)
	}
	return days
}
