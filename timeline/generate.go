package timeline

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
)

// Generate sample out for different file types, support csv, json and toml
func GenerateSample(filetype string, w io.Writer) error {
	var (
		tasks []SampleInput   // For csv
		all   SampleJsonInput // For json
	)
	filetype = strings.ToLower(filetype)

	// Sample data
	tasks = []SampleInput{
		SampleInput{
			Seq:   "1",
			Title: "EDA",
			Start: "2021-02-04",
			End:   "2021-02-10",
		},
		SampleInput{
			Seq:   "2",
			Title: "Build Model",
			Start: "2021-02-10",
			End:   "2021-03-10",
		},
		SampleInput{
			Seq:   "2",
			Title: "Build Super Model",
			Start: "2021-02-20",
			End:   "2021-02-23",
		},
		SampleInput{
			Seq:   "3",
			Title: "Evaluate Model",
			Start: "2021-03-05",
			End:   "2021-03-20",
		},
	}
	all.Tasks = tasks

	writer := bufio.NewWriter(w)

	// save to different format
	switch filetype {

	case "csv":

		content, err := gocsv.MarshalString(&tasks)
		if err != nil {
			return err
		}

		_, err = writer.WriteString(content)
		if err != nil {
			return err
		}
		err = writer.Flush()
		if err != nil {
			return err
		}

	case "json":
		_, err := writer.WriteString(PrettyPrint(all))
		if err != nil {
			return err
		}
		err = writer.Flush()
		if err != nil {
			return err
		}

	case "toml":
		return fmt.Errorf("File type is not supported - %s \n", filetype)

	default:
		return fmt.Errorf("File type is not supported - %s \n", filetype)
	}

	return nil
}
