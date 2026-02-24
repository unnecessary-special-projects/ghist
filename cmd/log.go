package cmd

import (
	"fmt"
	"os"

	"github.com/unnecessary-special-projects/ghist/internal/project"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log [message]",
	Short: "Record a decision or note to the event timeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		typ, _ := cmd.Flags().GetString("type")
		taskID, _ := cmd.Flags().GetInt64("task")

		var taskIDPtr *int64
		if cmd.Flags().Changed("task") {
			taskIDPtr = &taskID
		}

		event, err := s.CreateEvent(typ, args[0], "{}", taskIDPtr)
		if err != nil {
			return err
		}

		if err := project.UpdateContext(root, s); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to update context: %v\n", err)
		}

		fmt.Printf("Logged event #%d: %s\n", event.ID, event.Message)
		return nil
	},
}

func init() {
	logCmd.Flags().StringP("type", "t", "log", "Event type (log, decision, note)")
	logCmd.Flags().Int64P("task", "T", 0, "Link to task ID")
	rootCmd.AddCommand(logCmd)
}
