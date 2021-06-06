package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/billyla/timeline/timeline"
)

var (
	input  string
	output string
)

// createCmd creates the timeline html file
var createCmd = &cobra.Command{
	Use:     "create [--input sample.csv] [--ouput output.html]",
	Short:   "create a timeline.",
	Long:    `create a timeline.`,
	Aliases: []string{"c"},
	Example: `
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("")
		fmt.Printf("Input - %s \n", input)
		fmt.Printf("Output - %s \n", output)
		fmt.Println("")

		// Get file type
		filetype, err := getFileType(input)
		if err != nil {
			return err
		}

		// io.Reader
		f, err := os.Open(input)
		defer f.Close()
		if err != nil {
			return err
		}

		tasks, err := timeline.ParseFile(f, filetype)
		if err != nil {
			return err
		}

		// Write output
		f, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)

		err = timeline.CreateTimeline(w, tasks)
		if err != nil {
			return (err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&input, "input", "i", "sample.csv", "Input file")
	createCmd.Flags().StringVarP(&output, "output", "o", "output.html", "Output file")
}

// Get file type from the input, result should be csv, json or toml
func getFileType(filename string) (string, error) {
	splits := strings.Split(filename, ".")
	if len(splits) != 2 {
		return "", fmt.Errorf("Invalid file format, please check. Input - %s \n", filename)
	}
	filetype := splits[1]
	switch filetype {
	case "csv", "json", "toml":
		// pass

	default:
		return "", fmt.Errorf("Invalid file format, please check. Input - %s \n", filename)
	}
	return filetype, nil
}
