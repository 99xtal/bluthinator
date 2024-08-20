package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// programCmd represents the program command
var programCmd = &cobra.Command{
	Use:   "program",
	Short: "Run Gob's Program",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Gob's Program: (Y/N)?:\n? ")
			input, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
	
			input = strings.TrimSpace(input)
			if input == "y" || input == "Y" {
				for {
					for i := 0; i < 5; i++ {
						fmt.Print("Penus ")
					}
					fmt.Print("\n")
				}
			} else if input == "n" || input == "N" {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(programCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// programCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// programCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
