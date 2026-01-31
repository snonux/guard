package manager

import (
	"fmt"
	"os"
	"strconv"
)

// ShowConfig displays the current configuration from the registry
func (m *Manager) ShowConfig() error {
	if m.security == nil {
		return fmt.Errorf(".guardfile not found. Run 'guard init' first")
	}

	mode := m.security.GetDefaultFileMode()
	owner := m.security.GetDefaultFileOwner()
	group := m.security.GetDefaultFileGroup()

	// Format output per CLI-INTERFACE-SPECS.md
	fmt.Println("Configuration:")
	fmt.Printf("  Mode:  %04o\n", mode.Perm())
	fmt.Printf("  Owner: %s\n", formatConfigValue(owner))
	fmt.Printf("  Group: %s\n", formatConfigValue(group))

	return nil
}

// SetConfig updates guard configuration with one or more values
// Parameters with non-nil pointers are updated, nil means "don't change"
func (m *Manager) SetConfig(modeStr *string, owner *string, group *string) error {
	if m.security == nil {
		return fmt.Errorf(".guardfile not found. Run 'guard init' first")
	}

	// Check that at least one parameter is provided
	if modeStr == nil && owner == nil && group == nil {
		return fmt.Errorf("no configuration values provided")
	}

	// Track what we're updating
	var updates []string

	// Check if any files/collections are guarded (warning only)
	m.checkAndWarnGuardedFiles()

	// Update mode if provided
	if modeStr != nil {
		mode, err := parseOctalMode(*modeStr)
		if err != nil {
			return fmt.Errorf("invalid mode: %w", err)
		}

		if err := m.security.SetDefaultFileMode(mode); err != nil {
			return fmt.Errorf("failed to set mode: %w", err)
		}

		updates = append(updates, fmt.Sprintf("Mode:  %04o", mode.Perm()))
	}

	// Update owner if provided (can be empty string to clear)
	if owner != nil {
		m.security.SetDefaultFileOwner(*owner)
		if *owner == "" {
			updates = append(updates, "Owner: (cleared)")
		} else {
			updates = append(updates, fmt.Sprintf("Owner: %s", *owner))
		}
	}

	// Update group if provided (can be empty string to clear)
	if group != nil {
		m.security.SetDefaultFileGroup(*group)
		if *group == "" {
			updates = append(updates, "Group: (cleared)")
		} else {
			updates = append(updates, fmt.Sprintf("Group: %s", *group))
		}
	}

	// Save registry
	if err := m.security.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Print what was updated
	if len(updates) > 0 {
		fmt.Println("Config updated:")
		for _, update := range updates {
			fmt.Printf("  %s\n", update)
		}
	}

	return nil
}

// SetConfigMode updates guard_mode configuration
func (m *Manager) SetConfigMode(modeStr string) error {
	if m.security == nil {
		return fmt.Errorf(".guardfile not found. Run 'guard init' first")
	}

	// Parse octal string to os.FileMode
	mode, err := parseOctalMode(modeStr)
	if err != nil {
		return fmt.Errorf("invalid mode: %w", err)
	}

	// Check if any files/collections are guarded (warning only)
	m.checkAndWarnGuardedFiles()

	// Set the mode (this validates)
	if err := m.security.SetDefaultFileMode(mode); err != nil {
		return fmt.Errorf("failed to set mode: %w", err)
	}

	// Save registry
	if err := m.security.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Config updated:")
	fmt.Printf("  Mode: %04o\n", mode.Perm())
	return nil
}

// SetConfigOwner updates guard_owner configuration
func (m *Manager) SetConfigOwner(owner string) error {
	if m.security == nil {
		return fmt.Errorf(".guardfile not found. Run 'guard init' first")
	}

	// Check if any files/collections are guarded (warning only)
	m.checkAndWarnGuardedFiles()

	// Set the owner (trims whitespace)
	m.security.SetDefaultFileOwner(owner)

	// Save registry
	if err := m.security.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Config updated:")
	if owner == "" {
		fmt.Println("  Owner: (cleared)")
	} else {
		fmt.Printf("  Owner: %s\n", owner)
	}
	return nil
}

// SetConfigGroup updates guard_group configuration
func (m *Manager) SetConfigGroup(group string) error {
	if m.security == nil {
		return fmt.Errorf(".guardfile not found. Run 'guard init' first")
	}

	// Check if any files/collections are guarded (warning only)
	m.checkAndWarnGuardedFiles()

	// Set the group (trims whitespace)
	m.security.SetDefaultFileGroup(group)

	// Save registry
	if err := m.security.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Config updated:")
	if group == "" {
		fmt.Println("  Group: (cleared)")
	} else {
		fmt.Printf("  Group: %s\n", group)
	}
	return nil
}

// checkAndWarnGuardedFiles checks if any files/collections are guarded and adds a warning
func (m *Manager) checkAndWarnGuardedFiles() {
	guardedFileCount := 0
	guardedCollCount := 0

	// Count guarded files
	for _, file := range m.security.GetRegisteredFiles() {
		isGuarded, err := m.security.GetRegisteredFileGuard(file)
		if err == nil && isGuarded {
			guardedFileCount++
		}
	}

	// Count guarded collections
	for _, coll := range m.security.GetRegisteredCollections() {
		isGuarded, err := m.security.GetRegisteredCollectionGuard(coll)
		if err == nil && isGuarded {
			guardedCollCount++
		}
	}

	if guardedFileCount > 0 || guardedCollCount > 0 {
		// Format per spec: multi-line warning message
		warning := NewWarning(
			WarningGeneric,
			fmt.Sprintf("%d file(s) and %d collection(s) are currently guarded.\nThe new config will only apply to future guard operations.\nTo apply the new config to existing guards, disable and re-enable them.", guardedFileCount, guardedCollCount),
		)
		m.AddWarning(warning)
	}
}

// formatConfigValue formats a config value for display
func formatConfigValue(value string) string {
	if value == "" {
		return "(empty)"
	}
	return value
}

// parseOctalMode parses an octal mode string and returns os.FileMode
func parseOctalMode(modeStr string) (os.FileMode, error) {
	// Parse as uint32 in base 8
	modeInt, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("not a valid octal number: %s", modeStr)
	}

	// Check range (000-777)
	if modeInt > 0777 {
		return 0, fmt.Errorf("mode must be between 000 and 777, got: %s", modeStr)
	}

	return os.FileMode(modeInt), nil
}
