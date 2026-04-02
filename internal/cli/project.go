package cli

import (
	"fmt"
	"os"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "Add, list, or remove projects for skill installation",
}

func init() {
	rootCmd.AddCommand(projectCmd)
}

// projectAddCmd adds a new project
var projectAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a project",
	Run:   runProjectAdd,
}

func init() {
	projectCmd.AddCommand(projectAddCmd)
	projectAddCmd.Flags().StringVar(&flagName, "name", "", "Project name (required)")
	projectAddCmd.Flags().StringVar(&flagPath, "path", "", "Project path (required, defaults to current directory)")
	projectAddCmd.MarkFlagRequired("name")
}

func runProjectAdd(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	name := flagName
	path := flagPath

	// Default path to current directory
	if path == "" {
		cwd, err := getCurrentDir()
		if err != nil {
			Failf("Failed to get current directory: %v", err)
			return
		}
		path = cwd
	}

	// Validate path
	valid, err := h.ValidateProjectPath(path)
	if err != nil || !valid {
		Fail("Invalid project path: " + path)
		return
	}

	proj, err := h.CreateProject(name, path)
	if err != nil {
		Failf("Failed to create project: %v", err)
		return
	}

	if proj == nil {
		Fail("Project already exists: " + name)
		return
	}

	Success(map[string]interface{}{
		"project": proj,
		"message": fmt.Sprintf("Project '%s' added successfully", name),
	})
}

// projectListCmd lists all projects
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run:   runProjectList,
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}

func runProjectList(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	projects, err := h.ListProjects()
	if err != nil {
		Failf("Failed to list projects: %v", err)
		return
	}

	Success(map[string]interface{}{
		"projects": projects,
		"count":    len(projects),
	})
}

// projectRemoveCmd removes a project
var projectRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a project",
	Run:   runProjectRemove,
}

func init() {
	projectCmd.AddCommand(projectRemoveCmd)
	projectRemoveCmd.Flags().StringVar(&flagName, "name", "", "Project name (required)")
	projectRemoveCmd.MarkFlagRequired("name")
}

func runProjectRemove(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	name := flagName

	// Check if project exists
	exists, err := h.ProjectExists(name)
	if err != nil {
		Failf("Failed to check project: %v", err)
		return
	}
	if !exists {
		Fail("Project not found: " + name)
		return
	}

	if err := h.DeleteProject(name); err != nil {
		Failf("Failed to delete project: %v", err)
		return
	}

	Success(map[string]interface{}{
		"message": fmt.Sprintf("Project '%s' removed successfully", name),
	})
}

func getCurrentDir() (string, error) {
	return os.Getwd()
}
