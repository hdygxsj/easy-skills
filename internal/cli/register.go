package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a package with components to Hub",
	Long: `Register a package and its components to the local Hub.

Each --component flag specifies a component in the format "type:name:path".
  type: skill, agent, hook, or rule
  name: component name
  path: installed file/directory path (relative or absolute)

Example:
  easy-skills register --name my-pkg --target qoder \
    --component "skill:my-skill:.qoder/skills/my-skill" \
    --component "hook:my-hook:.qoder/hooks/my-hook"`,
	Run: runRegister,
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	registerCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	registerCmd.Flags().StringSliceVar(&flagComponents, "component", nil, `Component in format "type:name:path" (repeatable)`)
	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("target")
}

// parseComponent parses "type:name:path" format
func parseComponent(s string) (cType, name, path string, err error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid component format %q, expected type:name:path", s)
	}
	cType = strings.TrimSpace(parts[0])
	name = strings.TrimSpace(parts[1])
	path = strings.TrimSpace(parts[2])

	validTypes := map[string]bool{"skill": true, "agent": true, "hook": true, "rule": true}
	if !validTypes[cType] {
		return "", "", "", fmt.Errorf("invalid component type %q, must be one of: skill, agent, hook, rule", cType)
	}
	if name == "" || path == "" {
		return "", "", "", fmt.Errorf("component name and path cannot be empty")
	}
	return cType, name, path, nil
}

func runRegister(cmd *cobra.Command, args []string) {
	if len(flagComponents) == 0 {
		Fail("At least one --component is required")
		return
	}

	// Validate all components first
	type compDef struct {
		cType, name, path string
	}
	var defs []compDef
	for _, c := range flagComponents {
		cType, name, path, err := parseComponent(c)
		if err != nil {
			Failf("Invalid component: %v", err)
			return
		}
		defs = append(defs, compDef{cType, name, path})
	}

	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	// Check if package already exists
	exists, err := h.PackageExists(flagName, flagTarget)
	if err != nil {
		Failf("Failed to check package: %v", err)
		return
	}
	if exists {
		Fail("Package already exists. Use 'upgrade' to update it.")
		return
	}

	// Create package (no source - components hold the paths)
	pkg, err := h.CreatePackage(flagName, flagTarget, "", "{}")
	if err != nil {
		Failf("Failed to create package: %v", err)
		return
	}

	// Create initial version
	version, err := h.CreateVersion(pkg.ID, timestampVersion(), "{}")
	if err != nil {
		Failf("Failed to create version: %v", err)
		return
	}

	// Create components
	var components []map[string]string
	for _, d := range defs {
		comp, err := h.CreateComponent(pkg.ID, version.ID, d.name, d.cType, d.path)
		if err != nil {
			Failf("Failed to create component %s: %v", d.name, err)
			return
		}
		components = append(components, map[string]string{
			"id":   comp.ID,
			"name": comp.Name,
			"type": comp.Type,
			"path": comp.Path,
		})
	}

	Success(map[string]interface{}{
		"package":    pkg,
		"version":    version,
		"components": components,
		"message":    fmt.Sprintf("Package registered with %d component(s)", len(components)),
	})
}

func timestampVersion() int64 {
	return time.Now().Unix()
}
