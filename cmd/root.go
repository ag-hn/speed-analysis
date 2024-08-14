package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ag-hn/speed-analysis/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:     "sa",
	Short:   "Sa is a simple, configurable, and speed data analysis toolset",
	Version: "0.0.0",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enableLogging, err := cmd.Flags().GetBool("enable-logging")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("logging enabled")

		// If logging is enabled, logs will be output to debug.log.
		if enableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
		}

		cfg := tui.Config{
			EnableLogging:  enableLogging,
		}

		m := tui.New(cfg)

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			log.Fatal("Failed to start fm", err)
			os.Exit(1)
		}
	},
}

// Execute runs the root command and starts the application.
func Execute() {
	rootCmd.PersistentFlags().Bool("enable-logging", false, "Enable logging for FM")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
