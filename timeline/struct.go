package timeline

import (
	"html/template"
	"time"
)

type Output struct {
	Done     []Task
	Ongoing  []Task
	Timeline []Task
}

type AllTasks struct {
	Tasks []Task `json:"Tasks"`
}

type Task struct {
	Seq   string      `json:"Seq", csv:"Seq"`
	Title string      `json:"Title", csv:"Title"`
	Start template.JS `json:"Start", csv:"Start"` // yyyy-mm-dd format
	End   template.JS `json:"End", csv:"End"`     // yyyy-mm-dd format

	// Private
	StartD   time.Time
	EndD     time.Time
	Duration int // In number of business day.
}

// SampleJsonInput is similar to AllTasks but without private fields
type SampleJsonInput struct {
	Tasks []SampleInput `json:"Tasks"`
}

// SampleInput is similar to Task but without private fields
type SampleInput struct {
	Seq   string      `json:"Seq", csv:"Seq"`
	Title string      `json:"Title", csv:"Title"`
	Start template.JS `json:"Start", csv:"Start"` // yyyy-mm-dd format
	End   template.JS `json:"End", csv:"End"`     // yyyy-mm-dd format
}
