package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/billyla/timeline/timeline"
)

// generateCmd generates an input template in csv/json/toml format
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate sample input file in different formats.",
	Long:    `Generate sample input file in different formats.`,
	Aliases: []string{"g"},
	Example: `
  timeline generate csv
  timeline g csv
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Handle file type
		var filetype string
		if len(args) == 0 {
			filetype = "csv"
		} else if len(args) == 1 {
			filetype = strings.ToLower(args[0])
		} else {
			return fmt.Errorf("Please input one of the following formats. csv|json|toml\n")
		}

		filename := fmt.Sprintf("sample.%s", filetype) // sample.csv, sample.json, etc..
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)

		// Output
		err = timeline.GenerateSample(filetype, w)
		if err != nil {
			return err
		}
		fmt.Printf("Finished writing to output file - %s. \n", filename)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
