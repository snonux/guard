package manager

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileInfo contains display information for a registered file.
// Used by ShowFiles to return data instead of printing directly.
type FileInfo struct {
	Path        string
	Guard       bool
	Collections []string
}

// AddFiles registers files in the registry if they don't already exist.
// Per Requirement 2.3: Warns if files are missing on disk, but continues processing.
// Per Requirement 2.4: Idempotent - ignores files already registered (no warning).
func (m *Manager) AddFiles(paths []string) error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	if len(paths) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Validate all paths first (security check happens regardless of file existence)
	if err := m.security.ValidatePaths(paths); err != nil {
		return err
	}

	// Check which files exist
	existing, missing := m.fs.CheckFilesExist(paths)

	// Warn about missing files (context: files were not registered)
	if len(missing) > 0 {
		m.AddWarning(NewWarning(WarningFileMissing, "not_registered", missing...))
	}

	// Register existing files
	for _, path := range existing {
		// Check if already registered (idempotent)
		if m.security.IsRegisteredFile(path) {
			// Silent per Requirement 2.4
			continue
		}

		// Get current file info
		mode, owner, group, err := m.fs.GetFileInfo(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get file info for %s: %v", path, err))
			continue
		}

		// Register file
		if err := m.security.RegisterFile(path, mode, owner, group); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to register %s: %v", path, err))
			continue
		}
	}

	// Save registry
	if err := m.SaveRegistry(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	return nil
}

// RemoveFiles removes files from the registry and restores their original permissions.
// Per CLI-INTERFACE-SPECS.md lines 46-48, operation order:
//  1. Remove file from all collections
//  2. Restore original permissions (error if fails)
//  3. Remove from registry
//
// Per Requirement 2.6: Warns if files are not in registry.
// Per Requirement 11.6: Lenient - allows removal even if file deleted outside guard (with warning).
func (m *Manager) RemoveFiles(paths []string) error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	if len(paths) == 0 {
		return fmt.Errorf("no files specified")
	}

	for _, path := range paths {
		// Check if file is registered
		if !m.security.IsRegisteredFile(path) {
			m.AddWarning(NewWarning(WarningFileNotInRegistry, "", path))
			continue
		}

		// Step 1: Remove from all collections
		m.security.RemoveRegisteredFileFromAllRegisteredCollections(path)

		// Step 2: Restore original permissions if file is guarded
		// Get original metadata first to check if restoration is needed
		owner, group, mode, guard, err := m.security.GetRegisteredFileConfig(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get config for %s: %v", path, err))
			continue
		}

		// Only restore if currently guarded
		if guard {
			// Clear immutable flag first (must be done before chmod)
			if err := m.fs.ClearImmutable(path); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to clear immutable flag for %s: %v", path, err))
				continue
			}

			if err := m.fs.RestorePermissions(path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to restore permissions for %s: %v", path, err))
				continue
			}
		} else if !m.fs.FileExists(path) {
			// File not guarded and missing on disk - warn but continue
			m.AddWarning(NewWarning(WarningFileMissing, "", path))
		}

		// Step 3: Remove from registry
		if err := m.security.UnregisterFile(path, false); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to unregister %s: %v", path, err))
			continue
		}
	}

	// Save registry
	if err := m.SaveRegistry(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	return nil
}

// ToggleFiles toggles the guard status of files.
// Per Requirement 2.7: Adds missing files to registry first, then toggles.
func (m *Manager) ToggleFiles(paths []string) error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	if len(paths) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Check which files exist
	existing, missing := m.fs.CheckFilesExist(paths)

	// Warn about missing files
	if len(missing) > 0 {
		m.AddWarning(NewWarning(WarningFileMissing, "", missing...))
	}

	// Track files to toggle and their new guard states
	type fileToggle struct {
		path     string
		newGuard bool
	}
	var toggles []fileToggle

	// Phase 1: Register files and update registry state (in memory)
	for _, path := range existing {
		// Add to registry if not present
		if !m.security.IsRegisteredFile(path) {
			mode, owner, group, err := m.fs.GetFileInfo(path)
			if err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to get file info for %s: %v", path, err))
				continue
			}

			if err := m.security.RegisterFile(path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to register %s: %v", path, err))
				continue
			}
		}

		// Get current guard status
		guard, err := m.security.GetRegisteredFileGuard(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get guard status for %s: %v", path, err))
			continue
		}

		// Toggle guard status in registry
		newGuard := !guard
		if err := m.security.SetRegisteredFileGuard(path, newGuard); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set guard flag for %s: %v", path, err))
			continue
		}

		toggles = append(toggles, fileToggle{path: path, newGuard: newGuard})
	}

	// Phase 2: Save registry BEFORE applying filesystem permissions
	// This ensures the registry is saved while it's still writable
	if err := m.SaveRegistry(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	// Phase 3: Apply filesystem permissions
	for _, toggle := range toggles {
		if toggle.newGuard {
			// Enabling guard: apply guard permissions, then set immutable
			mode := m.security.GetDefaultFileMode()
			owner := m.security.GetDefaultFileOwner()
			group := m.security.GetDefaultFileGroup()

			if err := m.fs.ApplyPermissions(toggle.path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to apply guard permissions to %s: %v", toggle.path, err))
				continue
			}

			// Set immutable flag (auto-skips if not root)
			if err := m.fs.SetImmutable(toggle.path); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to set immutable flag for %s: %v", toggle.path, err))
			}
		} else {
			// Disabling guard: clear immutable first, then restore permissions
			if err := m.fs.ClearImmutable(toggle.path); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to clear immutable flag for %s: %v", toggle.path, err))
				continue
			}

			owner, group, mode, _, err := m.security.GetRegisteredFileConfig(toggle.path)
			if err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to get original config for %s: %v", toggle.path, err))
				continue
			}

			if err := m.fs.RestorePermissions(toggle.path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to restore permissions for %s: %v", toggle.path, err))
				continue
			}
		}
	}

	return nil
}

// EnableFiles enables guard protection on files.
// Per Requirement 5.1: Registers files if not in registry (with guard=false), then enables.
// Per Requirement 2b: Warns if files are missing on disk, does NOT change guard flag for missing files.
func (m *Manager) EnableFiles(paths []string) error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	if len(paths) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Check which files exist
	existing, missing := m.fs.CheckFilesExist(paths)

	// Warn about missing files
	if len(missing) > 0 {
		m.AddWarning(NewWarning(WarningFileMissing, "", missing...))
	}

	// Track files to enable
	var filesToEnable []string

	// Phase 1: Register files and update registry state (in memory)
	for _, path := range existing {
		// Add to registry if not present (Requirement 5.1)
		if !m.security.IsRegisteredFile(path) {
			mode, owner, group, err := m.fs.GetFileInfo(path)
			if err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to get file info for %s: %v", path, err))
				continue
			}

			// Register with guard=false initially
			if err := m.security.RegisterFile(path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to register %s: %v", path, err))
				continue
			}
		}

		// Set guard flag in registry
		if err := m.security.SetRegisteredFileGuard(path, true); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set guard flag for %s: %v", path, err))
			continue
		}

		filesToEnable = append(filesToEnable, path)
	}

	// Phase 2: Save registry BEFORE applying filesystem permissions
	if err := m.SaveRegistry(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	// Phase 3: Apply filesystem permissions
	for _, path := range filesToEnable {
		mode := m.security.GetDefaultFileMode()
		owner := m.security.GetDefaultFileOwner()
		group := m.security.GetDefaultFileGroup()

		if err := m.fs.ApplyPermissions(path, mode, owner, group); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to apply guard permissions to %s: %v", path, err))
			continue
		}

		// Set immutable flag (auto-skips if not root)
		if err := m.fs.SetImmutable(path); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set immutable flag for %s: %v", path, err))
		}
	}

	return nil
}

// DisableFiles disables guard protection on files and restores original permissions.
// Per Requirement 5.2: Warns if files are missing on disk or not in registry.
func (m *Manager) DisableFiles(paths []string) error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	if len(paths) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Check which files exist
	existing, missing := m.fs.CheckFilesExist(paths)

	// Warn about missing files
	if len(missing) > 0 {
		m.AddWarning(NewWarning(WarningFileMissing, "", missing...))
	}

	// Process existing files
	for _, path := range existing {
		// Check if registered
		if !m.security.IsRegisteredFile(path) {
			m.AddWarning(NewWarning(WarningFileNotInRegistry, "", path))
			continue
		}

		// Get original permissions
		owner, group, mode, guard, err := m.security.GetRegisteredFileConfig(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get config for %s: %v", path, err))
			continue
		}

		// Only restore if currently guarded
		if guard {
			// Clear immutable flag first (must be done before chmod)
			if err := m.fs.ClearImmutable(path); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to clear immutable flag for %s: %v", path, err))
				continue
			}

			if err := m.fs.RestorePermissions(path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to restore permissions for %s: %v", path, err))
				continue
			}
		}

		// Set guard flag to false
		if err := m.security.SetRegisteredFileGuard(path, false); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set guard flag for %s: %v", path, err))
			continue
		}
	}

	// Save registry
	if err := m.SaveRegistry(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	return nil
}

// ShowFiles displays the guard status of files and their collection membership.
// If no paths are specified, shows all registered files.
func (m *Manager) ShowFiles(paths []string) ([]FileInfo, error) {
	if m.security == nil {
		return nil, fmt.Errorf("registry not loaded")
	}

	// Initialize slice to collect file information
	fileInfos := []FileInfo{}

	// If no paths specified, show all registered files
	if len(paths) == 0 {
		paths = m.security.GetRegisteredFiles()
	}

	for _, path := range paths {
		if !m.security.IsRegisteredFile(path) {
			m.AddWarning(NewWarning(WarningFileNotInRegistry, "", path))
			continue
		}

		// Convert path to absolute for consistent comparison
		// (GetRegisteredCollectionFiles returns absolute paths)
		absPath, err := filepath.Abs(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to resolve path %s: %v", path, err))
			continue
		}

		// Get guard status
		guard, err := m.security.GetRegisteredFileGuard(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get guard status for %s: %v", path, err))
			continue
		}

		// Get collections this file belongs to
		collections := m.security.GetRegisteredCollections()
		memberOf := []string{}
		for _, coll := range collections {
			files, err := m.security.GetRegisteredCollectionFiles(coll)
			if err != nil {
				continue
			}
			for _, f := range files {
				if f == absPath {
					memberOf = append(memberOf, coll)
					break
				}
			}
		}

		// Collect file information
		fileInfos = append(fileInfos, FileInfo{
			Path:        path,
			Guard:       guard,
			Collections: memberOf,
		})

		// Check if file exists on disk
		if !m.fs.FileExists(path) {
			m.AddWarning(NewWarning(WarningFileMissing, "", path))
		}
	}

	return fileInfos, nil
}

// Cleanup removes empty collections and non-existent files from the registry.
// Per Requirement 8.1.
// CleanupResult contains the results of a cleanup operation.
type CleanupResult struct {
	FilesRemoved       int
	CollectionsRemoved int
}

func (m *Manager) Cleanup() (*CleanupResult, error) {
	if m.security == nil {
		return nil, fmt.Errorf("registry not loaded")
	}

	result := &CleanupResult{}

	// Remove non-existent files
	files := m.security.GetRegisteredFiles()
	for _, path := range files {
		if !m.fs.FileExists(path) {
			// Unregister file (also removes from collections)
			if err := m.security.UnregisterFile(path, true); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to unregister %s: %v", path, err))
			} else {
				result.FilesRemoved++
			}
		}
	}

	// Remove empty collections
	collections := m.security.GetRegisteredCollections()
	for _, coll := range collections {
		files, err := m.security.GetRegisteredCollectionFiles(coll)
		if err != nil {
			continue
		}
		if len(files) == 0 {
			if err := m.security.UnregisterCollection(coll, true); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to unregister collection %s: %v", coll, err))
			} else {
				result.CollectionsRemoved++
			}
		}
	}

	// Save registry
	if err := m.SaveRegistry(); err != nil {
		return nil, fmt.Errorf("failed to save registry: %w", err)
	}

	return result, nil
}

// ResetResult contains the results of a reset operation.
type ResetResult struct {
	FilesDisabled       int
	CollectionsDisabled int
}

// Reset disables guard for all files and collections.
// Per Requirement 8.2: Warns for missing files and recommends cleanup.
func (m *Manager) Reset() (*ResetResult, error) {
	if m.security == nil {
		return nil, fmt.Errorf("registry not loaded")
	}

	result := &ResetResult{}

	// Disable guard for all files
	files := m.security.GetRegisteredFiles()
	for _, path := range files {
		// Get config
		owner, group, mode, guard, err := m.security.GetRegisteredFileConfig(path)
		if err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to get config for %s: %v", path, err))
			continue
		}

		// Restore permissions if guarded
		if guard {
			// Clear immutable flag first (must be done before chmod)
			if err := m.fs.ClearImmutable(path); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to clear immutable flag for %s: %v", path, err))
				continue
			}

			if err := m.fs.RestorePermissions(path, mode, owner, group); err != nil {
				m.AddError(fmt.Sprintf("Error: Failed to restore permissions for %s: %v", path, err))
				continue
			}
			result.FilesDisabled++
		} else if !m.fs.FileExists(path) {
			// File not guarded and missing - warn but continue
			m.AddWarning(NewWarning(WarningFileMissing, "", path))
		}

		// Set guard flag to false
		if err := m.security.SetRegisteredFileGuard(path, false); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set guard flag for %s: %v", path, err))
		}
	}

	// Disable guard for all collections
	collections := m.security.GetRegisteredCollections()
	for _, coll := range collections {
		// Get current guard state to count only disabled collections
		guard, err := m.security.GetRegisteredCollectionGuard(coll)
		if err == nil && guard {
			result.CollectionsDisabled++
		}

		if err := m.security.SetRegisteredCollectionGuard(coll, false); err != nil {
			m.AddError(fmt.Sprintf("Error: Failed to set guard flag for collection %s: %v", coll, err))
		}
	}

	// Save registry
	if err := m.SaveRegistry(); err != nil {
		return nil, fmt.Errorf("failed to save registry: %w", err)
	}

	return result, nil
}

// Destroy runs reset, cleanup, verifies all permissions restored, and deletes .guardfile.
// Per Requirement 8.3: Only deletes .guardfile if verification succeeds.
func (m *Manager) Destroy() error {
	if m.security == nil {
		return fmt.Errorf("registry not loaded")
	}

	// Step 1: Reset (disable all guards)
	resetResult, err := m.Reset()
	if err != nil {
		return fmt.Errorf("reset failed: %w", err)
	}

	// Print reset output per CLI-INTERFACE-SPECS.md
	fmt.Println("Reset complete:")
	if resetResult.FilesDisabled > 0 || resetResult.CollectionsDisabled > 0 {
		if resetResult.FilesDisabled > 0 {
			fmt.Printf("  Guard disabled for %d file(s)\n", resetResult.FilesDisabled)
		}
		if resetResult.CollectionsDisabled > 0 {
			fmt.Printf("  Guard disabled for %d collection(s)\n", resetResult.CollectionsDisabled)
		}
	} else {
		fmt.Println("  No guarded files or collections found")
	}

	// Check for reset errors and abort if any
	if m.HasErrors() {
		return fmt.Errorf("uninstall aborted. Fix errors and try again")
	}

	// Step 2: Cleanup (remove empty collections and missing files)
	cleanupResult, err := m.Cleanup()
	if err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	// Print cleanup output per CLI-INTERFACE-SPECS.md
	fmt.Println("Cleanup complete:")
	if cleanupResult.FilesRemoved > 0 || cleanupResult.CollectionsRemoved > 0 {
		fmt.Printf("  Removed %d file(s) (file not found)\n", cleanupResult.FilesRemoved)
		fmt.Printf("  Removed %d collection(s) (empty)\n", cleanupResult.CollectionsRemoved)
	} else {
		fmt.Println("  No stale entries found")
	}

	// Step 3: Verify all existing files have restored permissions
	files := m.security.GetRegisteredFiles()
	verificationFailed := false
	for _, path := range files {
		if !m.fs.FileExists(path) {
			continue // Skip missing files
		}

		// Get expected (original) permissions
		expectedOwner, expectedGroup, expectedMode, guard, err := m.security.GetRegisteredFileConfig(path)
		if err != nil {
			return fmt.Errorf("verification failed: cannot get config for %s: %w", path, err)
		}

		// Guard should be false after reset
		if guard {
			verificationFailed = true
			m.AddError(fmt.Sprintf("Error: File %s still has guard enabled", path))
		}

		// Verify actual permissions match expected
		actualMode, actualOwner, actualGroup, err := m.fs.GetFileInfo(path)
		if err != nil {
			return fmt.Errorf("verification failed: cannot get file info for %s: %w", path, err)
		}

		if actualMode != expectedMode || actualOwner != expectedOwner || actualGroup != expectedGroup {
			verificationFailed = true
			m.AddError(fmt.Sprintf("Error: File %s permissions not restored (expected: mode=%o owner=%s group=%s, got: mode=%o owner=%s group=%s)",
				path, expectedMode, expectedOwner, expectedGroup, actualMode, actualOwner, actualGroup))
		}
	}

	// Step 4: Only delete .guardfile if verification succeeded
	if verificationFailed || m.HasErrors() {
		return fmt.Errorf("destroy verification failed - .guardfile preserved. Fix errors and try again")
	}

	// Delete .guardfile
	if err := os.Remove(m.registryPath); err != nil {
		return fmt.Errorf("failed to delete .guardfile: %w", err)
	}
	fmt.Println("Removed .guardfile")
	fmt.Println("Uninstall complete")

	return nil
}
