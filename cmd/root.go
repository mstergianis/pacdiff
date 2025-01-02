/*
Copyright © 2024 Michael Stergianis michael.stergianis@gmail.com
*/
package cmd

import (
	"os"

	"github.com/mstergianis/pacdiff/pkg/depthparser"
	"github.com/mstergianis/pacdiff/pkg/differ"
	"github.com/mstergianis/pacdiff/pkg/printer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pacdiff LEFT-PACKAGE RIGHT-PACKAGE",
	Short: "Golang semantic package differ",
	Long: `pacdiff is a cli tool for diffing golang packages.

pacdiff <LEFT-PACKAGE> <RIGHT-PACKAGE>
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MatchAll(cobra.ExactArgs(2)),
	RunE: func(cmd *cobra.Command, args []string) error {
		d := differ.NewDiffer(differ.WithPackages(args[0], args[1]))
		diff, err := d.Diff()
		if err != nil {
			return err
		}

		depthStop, err := depthparser.Parse(cmd.Flags().Lookup("depth-delimiter").Value.String())
		if err != nil {
			return err
		}

		p := printer.NewPrinter(
			printer.WithDepth(depthStop),
		)
		p.Print(diff)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	rootCmd.Flags().StringP("depth-delimiter", "d", "2s", `provide a delimiter to use as a depth

grammar:
  depth-delimiter = { number } space-type .
  space-type = "s" | "t" .
  number = decimal-digit { decimal-digit } .
  decimal-digit = "0" | "1" | ... | "8" | "9" .
`)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
