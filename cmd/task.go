package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/unnecessary-special-projects/ghist/internal/models"
	"github.com/unnecessary-special-projects/ghist/internal/output"
	"github.com/unnecessary-special-projects/ghist/internal/project"
	"github.com/unnecessary-special-projects/ghist/internal/store"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

// --- task add ---

var taskAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		milestone, _ := cmd.Flags().GetString("milestone")
		priority, _ := cmd.Flags().GetString("priority")
		taskType, _ := cmd.Flags().GetString("type")
		legacyID, _ := cmd.Flags().GetString("legacy-id")

		task, err := s.CreateTask(store.CreateTaskInput{
			Title:       args[0],
			Description: description,
			Status:      status,
			Milestone:   milestone,
			Priority:    priority,
			Type:        taskType,
			LegacyID:    legacyID,
		})
		if err != nil {
			return err
		}

		if err := project.UpdateContext(root, s); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to update context: %v\n", err)
		}

		fmt.Printf("Created task %s: %s\n", task.RefID, task.Title)
		return nil
	},
}

// --- task list ---

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		status, _ := cmd.Flags().GetString("status")
		milestone, _ := cmd.Flags().GetString("milestone")
		priority, _ := cmd.Flags().GetString("priority")
		taskType, _ := cmd.Flags().GetString("type")
		asJSON, _ := cmd.Flags().GetBool("json")

		tasks, err := s.ListTasks(status, milestone, priority, taskType)
		if err != nil {
			return err
		}

		if asJSON {
			data, err := json.MarshalIndent(tasks, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		output.PrintTaskTable(tasks)
		return nil
	},
}

// --- task show ---

var taskShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		id, err := models.ParseTaskID(args[0])
		if err != nil {
			return err
		}

		task, err := s.GetTask(id)
		if err != nil {
			return err
		}

		events, err := s.ListEventsByTask(id)
		if err != nil {
			return err
		}

		output.PrintTaskDetail(task, events)
		return nil
	},
}

// --- task update ---

var taskUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		id, err := models.ParseTaskID(args[0])
		if err != nil {
			return err
		}

		u := store.TaskUpdate{}

		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			u.Title = &v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			u.Description = &v
		}
		if cmd.Flags().Changed("status") {
			v, _ := cmd.Flags().GetString("status")
			u.Status = &v
		}
		if cmd.Flags().Changed("milestone") {
			v, _ := cmd.Flags().GetString("milestone")
			u.Milestone = &v
		}
		if cmd.Flags().Changed("commit-hash") {
			v, _ := cmd.Flags().GetString("commit-hash")
			u.CommitHash = &v
		}
		if cmd.Flags().Changed("plan") {
			v, _ := cmd.Flags().GetString("plan")
			u.Plan = &v
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetString("priority")
			u.Priority = &v
		}
		if cmd.Flags().Changed("type") {
			v, _ := cmd.Flags().GetString("type")
			u.Type = &v
		}
		if cmd.Flags().Changed("legacy-id") {
			v, _ := cmd.Flags().GetString("legacy-id")
			u.LegacyID = &v
		}
		planStdin, _ := cmd.Flags().GetBool("plan-stdin")
		if planStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading plan from stdin: %w", err)
			}
			v := string(data)
			u.Plan = &v
		}

		task, err := s.UpdateTask(id, u)
		if err != nil {
			return err
		}

		if err := project.UpdateContext(root, s); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to update context: %v\n", err)
		}

		fmt.Printf("Updated task %s: %s [%s]\n", task.RefID, task.Title, task.Status)
		return nil
	},
}

// --- task delete ---

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		id, err := models.ParseTaskID(args[0])
		if err != nil {
			return err
		}

		if err := s.DeleteTask(id); err != nil {
			return err
		}

		if err := project.UpdateContext(root, s); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to update context: %v\n", err)
		}

		fmt.Printf("Deleted task #%d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskAddCmd.Flags().StringP("description", "d", "", "Task description")
	taskAddCmd.Flags().StringP("status", "s", "", "Task status (todo, in_planning, in_progress, done, blocked)")
	taskAddCmd.Flags().StringP("milestone", "m", "", "Milestone name")
	taskAddCmd.Flags().StringP("priority", "p", "", "Priority (low, medium, high, urgent)")
	taskAddCmd.Flags().StringP("type", "t", "", "Type (bug, feature, improvement, chore)")
	taskAddCmd.Flags().String("legacy-id", "", "Legacy ID from external system")
	taskCmd.AddCommand(taskAddCmd)

	taskListCmd.Flags().StringP("status", "s", "", "Filter by status")
	taskListCmd.Flags().StringP("milestone", "m", "", "Filter by milestone")
	taskListCmd.Flags().StringP("priority", "p", "", "Filter by priority")
	taskListCmd.Flags().StringP("type", "t", "", "Filter by type")
	taskListCmd.Flags().Bool("json", false, "Output as JSON")
	taskCmd.AddCommand(taskListCmd)

	taskCmd.AddCommand(taskShowCmd)

	taskUpdateCmd.Flags().String("title", "", "New title")
	taskUpdateCmd.Flags().StringP("description", "d", "", "New description")
	taskUpdateCmd.Flags().StringP("status", "s", "", "New status (todo, in_planning, in_progress, done, blocked)")
	taskUpdateCmd.Flags().StringP("milestone", "m", "", "New milestone")
	taskUpdateCmd.Flags().String("commit-hash", "", "Associated commit hash")
	taskUpdateCmd.Flags().String("plan", "", "Implementation plan text")
	taskUpdateCmd.Flags().Bool("plan-stdin", false, "Read plan from stdin")
	taskUpdateCmd.Flags().StringP("priority", "p", "", "Priority (low, medium, high, urgent)")
	taskUpdateCmd.Flags().StringP("type", "t", "", "Type (bug, feature, improvement, chore)")
	taskUpdateCmd.Flags().String("legacy-id", "", "Legacy ID from external system")
	taskCmd.AddCommand(taskUpdateCmd)

	taskCmd.AddCommand(taskDeleteCmd)
}

// openStore finds the project root and opens the store.
func openStore() (string, *store.Store, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil, fmt.Errorf("getting working directory: %w", err)
	}

	root, err := project.FindRoot(cwd)
	if err != nil {
		return "", nil, fmt.Errorf("not a ghist project (run 'ghist init' first): %w", err)
	}

	s, err := store.Open(project.GhistDirPath(root))
	if err != nil {
		return "", nil, fmt.Errorf("opening database: %w", err)
	}

	return root, s, nil
}
