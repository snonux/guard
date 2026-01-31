package registry

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// ============================================================================
// Test Category 2.1: Registry Creation & Initialization
// ============================================================================

func TestNewRegistryWithValidDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if reg == nil {
		t.Fatal("Registry is nil")
	}

	// Verify config was set correctly
	mode := reg.GetDefaultFileMode()
	if mode != 0640 {
		t.Errorf("Expected mode 0640, got %o", mode)
	}

	owner := reg.GetDefaultFileOwner()
	if owner != "root" {
		t.Errorf("Expected owner 'root', got '%s'", owner)
	}

	group := reg.GetDefaultFileGroup()
	if group != "wheel" {
		t.Errorf("Expected group 'wheel', got '%s'", group)
	}
}

func TestNewRegistryWithInvalidMode(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	testCases := []struct {
		name string
		mode string
	}{
		{"out of range", "0888"},
		{"invalid digits", "999"},
		{"too many digits", "07777"},
		{"non-numeric", "abc"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defaults := &RegistryDefaults{
				GuardMode:  tc.mode,
				GuardOwner: "root",
				GuardGroup: "wheel",
			}

			_, err := NewRegistry(registryPath, defaults, false)
			if err == nil {
				t.Errorf("Expected error for invalid mode '%s', got nil", tc.mode)
			}
		})
	}
}

func TestNewRegistryWithNilDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	_, err := NewRegistry(registryPath, nil, false)
	if err == nil {
		t.Error("Expected error for nil defaults, got nil")
	}
}

func TestNewRegistryWithValidModeRange(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		mode     string
		expected os.FileMode
	}{
		{"000", 0000},
		{"0000", 0000},
		{"644", 0644},
		{"0644", 0644},
		{"777", 0777},
		{"0777", 0777},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			registryPath := filepath.Join(tmpDir, ".guardfile_"+tc.mode)
			defaults := &RegistryDefaults{
				GuardMode:  tc.mode,
				GuardOwner: "user",
				GuardGroup: "group",
			}

			reg, err := NewRegistry(registryPath, defaults, false)
			if err != nil {
				t.Fatalf("NewRegistry failed for mode '%s': %v", tc.mode, err)
			}

			mode := reg.GetDefaultFileMode()
			if mode != tc.expected {
				t.Errorf("Expected mode %o, got %o", tc.expected, mode)
			}
		})
	}
}

func TestNewRegistryWithEmptyOwnerGroup(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "",
		GuardGroup: "",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed with empty owner/group: %v", err)
	}

	owner := reg.GetDefaultFileOwner()
	if owner != "" {
		t.Errorf("Expected empty owner, got '%s'", owner)
	}

	group := reg.GetDefaultFileGroup()
	if group != "" {
		t.Errorf("Expected empty group, got '%s'", group)
	}
}

func TestNewRegistryFileAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	// Create first registry
	reg1, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("First NewRegistry failed: %v", err)
	}

	// Save it to create the file
	if err := reg1.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Try to create second registry without overwrite
	_, err = NewRegistry(registryPath, defaults, false)
	if err == nil {
		t.Error("Expected error when file exists and overwrite=false, got nil")
	}

	// Try with overwrite=true
	reg2, err := NewRegistry(registryPath, defaults, true)
	if err != nil {
		t.Errorf("NewRegistry with overwrite=true failed: %v", err)
	}
	if reg2 == nil {
		t.Error("Registry is nil with overwrite=true")
	}
}

func TestLoadRegistryMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, "nonexistent.guardfile")

	_, err := LoadRegistry(registryPath)
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

// ============================================================================
// Test Category 2.2: YAML Persistence
// ============================================================================

func TestSaveAndLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}

	// Create and populate registry
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Add some files
	if err := reg.RegisterFile("./file1.txt", 0644, "user1", "group1"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile("./file2.txt", 0600, "user2", "group2"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Set guard flag on one file
	if err := reg.SetRegisteredFileGuard("./file1.txt", true); err != nil {
		t.Fatalf("SetRegisteredFileGuard failed: %v", err)
	}

	// Add collections
	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Add files to collection
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./file1.txt", "./file2.txt"}); err != nil {
		t.Fatalf("AddRegisteredFilesToRegisteredCollections failed: %v", err)
	}

	// Save registry
	if err := reg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load registry
	loaded, err := LoadRegistry(registryPath)
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	// Verify config
	if loaded.GetDefaultFileMode() != 0640 {
		t.Errorf("Loaded mode mismatch: expected 0640, got %o", loaded.GetDefaultFileMode())
	}
	if loaded.GetDefaultFileOwner() != "testuser" {
		t.Errorf("Loaded owner mismatch: expected 'testuser', got '%s'", loaded.GetDefaultFileOwner())
	}
	if loaded.GetDefaultFileGroup() != "testgroup" {
		t.Errorf("Loaded group mismatch: expected 'testgroup', got '%s'", loaded.GetDefaultFileGroup())
	}

	// Verify files
	if !loaded.IsRegisteredFile("./file1.txt") {
		t.Error("file1.txt not found in loaded registry")
	}
	if !loaded.IsRegisteredFile("./file2.txt") {
		t.Error("file2.txt not found in loaded registry")
	}

	// Verify file metadata
	mode1, err := loaded.GetRegisteredFileMode("./file1.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileMode failed: %v", err)
	}
	if mode1 != 0644 {
		t.Errorf("file1 mode mismatch: expected 0644, got %o", mode1)
	}

	guard1, err := loaded.GetRegisteredFileGuard("./file1.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileGuard failed: %v", err)
	}
	if !guard1 {
		t.Error("file1 guard flag should be true")
	}

	// Verify collections
	if !loaded.IsRegisteredCollection("coll1") {
		t.Error("coll1 not found in loaded registry")
	}

	files, err := loaded.GetRegisteredCollectionFiles("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 files in coll1, got %d", len(files))
	}
}

func TestLoadCorruptedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	// Write invalid YAML
	invalidYAML := "config: {{{"
	if err := os.WriteFile(registryPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write invalid YAML: %v", err)
	}

	_, err := LoadRegistry(registryPath)
	if err == nil {
		t.Error("Expected error for corrupted YAML, got nil")
	}
}

func TestLoadInvalidModeInYAML(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	// Write YAML with invalid mode
	invalidYAML := `config:
  guard_mode: "0888"
  guard_owner: "root"
  guard_group: "wheel"
files: []
collections: []
`
	if err := os.WriteFile(registryPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write YAML: %v", err)
	}

	_, err := LoadRegistry(registryPath)
	if err == nil {
		t.Error("Expected error for invalid mode in YAML, got nil")
	}
}

// ============================================================================
// Test Category 2.3: File Operations
// ============================================================================

func TestRegisterFile(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	err = reg.RegisterFile("./test.txt", 0644, "user", "staff")
	if err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	if !reg.IsRegisteredFile("./test.txt") {
		t.Error("File should be registered")
	}

	// Verify metadata
	mode, err := reg.GetRegisteredFileMode("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileMode failed: %v", err)
	}
	if mode != 0644 {
		t.Errorf("Expected mode 0644, got %o", mode)
	}

	owner, err := reg.GetRegisteredFileOwner("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileOwner failed: %v", err)
	}
	if owner != "user" {
		t.Errorf("Expected owner 'user', got '%s'", owner)
	}

	group, err := reg.GetRegisteredFileGroup("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileGroup failed: %v", err)
	}
	if group != "staff" {
		t.Errorf("Expected group 'staff', got '%s'", group)
	}

	guard, err := reg.GetRegisteredFileGuard("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileGuard failed: %v", err)
	}
	if guard {
		t.Error("New file should have guard=false")
	}
}

func TestRegisterFileDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register first time
	err = reg.RegisterFile("./test.txt", 0644, "user", "staff")
	if err != nil {
		t.Fatalf("First RegisterFile failed: %v", err)
	}

	// Register second time (should error)
	err = reg.RegisterFile("./test.txt", 0644, "user", "staff")
	if err == nil {
		t.Error("Expected error for duplicate registration, got nil")
	}
}

func TestUnregisterFile(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register file
	if err := reg.RegisterFile("./test.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Unregister it
	err = reg.UnregisterFile("./test.txt", false)
	if err != nil {
		t.Fatalf("UnregisterFile failed: %v", err)
	}

	if reg.IsRegisteredFile("./test.txt") {
		t.Error("File should not be registered after unregister")
	}
}

func TestUnregisterFileNotRegistered(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Try to unregister non-existent file with ignoreMissing=false
	err = reg.UnregisterFile("./nonexistent.txt", false)
	if err == nil {
		t.Error("Expected error for unregistering non-existent file, got nil")
	}

	// Try with ignoreMissing=true
	err = reg.UnregisterFile("./nonexistent.txt", true)
	if err != nil {
		t.Errorf("UnregisterFile with ignoreMissing=true should not error, got: %v", err)
	}
}

func TestUnregisterFileRemovesFromAllCollections(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register file
	if err := reg.RegisterFile("./test.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Create collections
	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}
	if err := reg.RegisterCollection("coll2", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Add file to both collections
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1", "coll2"}, []string{"./test.txt"}); err != nil {
		t.Fatalf("AddRegisteredFilesToRegisteredCollections failed: %v", err)
	}

	// Unregister file
	if err := reg.UnregisterFile("./test.txt", false); err != nil {
		t.Fatalf("UnregisterFile failed: %v", err)
	}

	// Verify file is removed from both collections
	files1, err := reg.GetRegisteredCollectionFiles("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files1) != 0 {
		t.Errorf("coll1 should have 0 files, got %d", len(files1))
	}

	files2, err := reg.GetRegisteredCollectionFiles("coll2")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files2) != 0 {
		t.Errorf("coll2 should have 0 files, got %d", len(files2))
	}
}

func TestGetRegisteredFiles(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Initially empty
	files := reg.GetRegisteredFiles()
	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}

	// Register files
	if err := reg.RegisterFile("./file1.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile("./file2.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	files = reg.GetRegisteredFiles()
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestSetAndGetFileMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterFile("./test.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Test SetRegisteredFileMode
	if err := reg.SetRegisteredFileMode("./test.txt", 0600); err != nil {
		t.Fatalf("SetRegisteredFileMode failed: %v", err)
	}
	mode, err := reg.GetRegisteredFileMode("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileMode failed: %v", err)
	}
	if mode != 0600 {
		t.Errorf("Expected mode 0600, got %o", mode)
	}

	// Test SetRegisteredFileGuard
	if err := reg.SetRegisteredFileGuard("./test.txt", true); err != nil {
		t.Fatalf("SetRegisteredFileGuard failed: %v", err)
	}
	guard, err := reg.GetRegisteredFileGuard("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileGuard failed: %v", err)
	}
	if !guard {
		t.Error("Guard should be true")
	}

	// Test SetRegisteredFileOwner
	prevOwner, err := reg.SetRegisteredFileOwner("./test.txt", "newuser")
	if err != nil {
		t.Fatalf("SetRegisteredFileOwner failed: %v", err)
	}
	if prevOwner != "user" {
		t.Errorf("Expected previous owner 'user', got '%s'", prevOwner)
	}
	owner, err := reg.GetRegisteredFileOwner("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileOwner failed: %v", err)
	}
	if owner != "newuser" {
		t.Errorf("Expected owner 'newuser', got '%s'", owner)
	}

	// Test SetRegisteredFileGroup
	prevGroup, err := reg.SetRegisteredFileGroup("./test.txt", "newgroup")
	if err != nil {
		t.Fatalf("SetRegisteredFileGroup failed: %v", err)
	}
	if prevGroup != "staff" {
		t.Errorf("Expected previous group 'staff', got '%s'", prevGroup)
	}
	group, err := reg.GetRegisteredFileGroup("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileGroup failed: %v", err)
	}
	if group != "newgroup" {
		t.Errorf("Expected group 'newgroup', got '%s'", group)
	}
}

func TestGetRegisteredFileConfig(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterFile("./test.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	if err := reg.SetRegisteredFileGuard("./test.txt", true); err != nil {
		t.Fatalf("SetRegisteredFileGuard failed: %v", err)
	}

	owner, group, mode, guard, err := reg.GetRegisteredFileConfig("./test.txt")
	if err != nil {
		t.Fatalf("GetRegisteredFileConfig failed: %v", err)
	}

	if owner != "user" {
		t.Errorf("Expected owner 'user', got '%s'", owner)
	}
	if group != "staff" {
		t.Errorf("Expected group 'staff', got '%s'", group)
	}
	if mode != 0644 {
		t.Errorf("Expected mode 0644, got %o", mode)
	}
	if !guard {
		t.Error("Guard should be true")
	}
}

// ============================================================================
// Test Category 2.4: Collection Operations
// ============================================================================

func TestRegisterCollection(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	err = reg.RegisterCollection("testcoll", []string{})
	if err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	if !reg.IsRegisteredCollection("testcoll") {
		t.Error("Collection should be registered")
	}

	// Verify guard flag is false for new collection
	guard, err := reg.GetRegisteredCollectionGuard("testcoll")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionGuard failed: %v", err)
	}
	if guard {
		t.Error("New collection should have guard=false")
	}
}

func TestRegisterCollectionDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register first time
	if err := reg.RegisterCollection("testcoll", []string{}); err != nil {
		t.Fatalf("First RegisterCollection failed: %v", err)
	}

	// Register second time (should error)
	err = reg.RegisterCollection("testcoll", []string{})
	if err == nil {
		t.Error("Expected error for duplicate collection, got nil")
	}
}

func TestUnregisterCollection(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterCollection("testcoll", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	err = reg.UnregisterCollection("testcoll", false)
	if err != nil {
		t.Fatalf("UnregisterCollection failed: %v", err)
	}

	if reg.IsRegisteredCollection("testcoll") {
		t.Error("Collection should not be registered after unregister")
	}
}

func TestUnregisterCollectionNotRegistered(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Try with ignoreMissing=false
	err = reg.UnregisterCollection("nonexistent", false)
	if err == nil {
		t.Error("Expected error for unregistering non-existent collection, got nil")
	}

	// Try with ignoreMissing=true
	err = reg.UnregisterCollection("nonexistent", true)
	if err != nil {
		t.Errorf("UnregisterCollection with ignoreMissing=true should not error, got: %v", err)
	}
}

func TestGetRegisteredCollections(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Initially empty
	colls := reg.GetRegisteredCollections()
	if len(colls) != 0 {
		t.Errorf("Expected 0 collections, got %d", len(colls))
	}

	// Register collections
	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}
	if err := reg.RegisterCollection("coll2", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	colls = reg.GetRegisteredCollections()
	if len(colls) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(colls))
	}
}

func TestCountFilesInCollection(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register files
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(tmpDir, "file3.txt")

	if err := reg.RegisterFile(file1, 0644, "user", "group"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile(file2, 0644, "user", "group"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile(file3, 0644, "user", "group"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Register collection with no files
	if err := reg.RegisterCollection("empty", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Test empty collection
	count, err := reg.CountFilesInCollection("empty")
	if err != nil {
		t.Fatalf("CountFilesInCollection failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 files in empty collection, got %d", count)
	}

	// Register collection with files
	if err := reg.RegisterCollection("mycoll", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Add files to collection
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"mycoll"}, []string{file1, file2, file3}); err != nil {
		t.Fatalf("AddRegisteredFilesToRegisteredCollections failed: %v", err)
	}

	// Test collection with 3 files
	count, err = reg.CountFilesInCollection("mycoll")
	if err != nil {
		t.Fatalf("CountFilesInCollection failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 files in mycoll, got %d", count)
	}

	// Test non-existent collection
	_, err = reg.CountFilesInCollection("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent collection, got nil")
	}
}

func TestAddRegisteredFilesToRegisteredCollections(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register files
	if err := reg.RegisterFile("./file1.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile("./file2.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Register collection
	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Add files to collection
	err = reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./file1.txt", "./file2.txt"})
	if err != nil {
		t.Fatalf("AddRegisteredFilesToRegisteredCollections failed: %v", err)
	}

	// Verify files are in collection
	files, err := reg.GetRegisteredCollectionFiles("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 files in collection, got %d", len(files))
	}
}

func TestAddRegisteredFilesToRegisteredCollectionsIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterFile("./file1.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Add file once
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./file1.txt"}); err != nil {
		t.Fatalf("First Add failed: %v", err)
	}

	// Add same file again (should be idempotent)
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./file1.txt"}); err != nil {
		t.Fatalf("Second Add failed: %v", err)
	}

	// Verify file appears only once
	files, err := reg.GetRegisteredCollectionFiles("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file (idempotent), got %d", len(files))
	}
}

func TestAddRegisteredFilesToRegisteredCollectionsNonExistentCollection(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterFile("./file1.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	// Try to add to non-existent collection (should error)
	err = reg.AddRegisteredFilesToRegisteredCollections([]string{"nonexistent"}, []string{"./file1.txt"})
	if err == nil {
		t.Error("Expected error for non-existent collection, got nil")
	}
}

func TestAddRegisteredFilesToRegisteredCollectionsNonRegisteredFile(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Try to add non-registered file (should error)
	err = reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./nonregistered.txt"})
	if err == nil {
		t.Error("Expected error for non-registered file, got nil")
	}
}

func TestRemoveRegisteredFilesFromRegisteredCollections(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Setup: register files and collection
	if err := reg.RegisterFile("./file1.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterFile("./file2.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}
	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}
	if err := reg.AddRegisteredFilesToRegisteredCollections([]string{"coll1"}, []string{"./file1.txt", "./file2.txt"}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Remove file1
	err = reg.RemoveRegisteredFilesFromRegisteredCollections([]string{"coll1"}, []string{"./file1.txt"})
	if err != nil {
		t.Fatalf("RemoveRegisteredFilesFromRegisteredCollections failed: %v", err)
	}

	// Verify only file2 remains
	files, err := reg.GetRegisteredCollectionFiles("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionFiles failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file remaining, got %d", len(files))
	}
	if len(files) > 0 && files[0] != "./file2.txt" {
		t.Errorf("Expected file2.txt, got %s", files[0])
	}
}

func TestSetAndGetCollectionGuard(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	if err := reg.RegisterCollection("coll1", []string{}); err != nil {
		t.Fatalf("RegisterCollection failed: %v", err)
	}

	// Initially false
	guard, err := reg.GetRegisteredCollectionGuard("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionGuard failed: %v", err)
	}
	if guard {
		t.Error("New collection should have guard=false")
	}

	// Set to true
	if err := reg.SetRegisteredCollectionGuard("coll1", true); err != nil {
		t.Fatalf("SetRegisteredCollectionGuard failed: %v", err)
	}

	guard, err = reg.GetRegisteredCollectionGuard("coll1")
	if err != nil {
		t.Fatalf("GetRegisteredCollectionGuard failed: %v", err)
	}
	if !guard {
		t.Error("Collection guard should be true")
	}
}

// ============================================================================
// Test Category 2.5: Thread-Safety Testing
// ============================================================================

func TestConcurrentRegisterFile(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Launch 10 goroutines registering different files
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			path := filepath.Join("./", "file"+string(rune('0'+idx))+".txt")
			if err := reg.RegisterFile(path, 0644, "user", "staff"); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent RegisterFile error: %v", err)
	}

	// Verify all files registered
	files := reg.GetRegisteredFiles()
	if len(files) != 10 {
		t.Errorf("Expected 10 files, got %d", len(files))
	}
}

func TestConcurrentSave(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register some files
	for i := 0; i < 5; i++ {
		path := filepath.Join("./", "file"+string(rune('0'+i))+".txt")
		if err := reg.RegisterFile(path, 0644, "user", "staff"); err != nil {
			t.Fatalf("RegisterFile failed: %v", err)
		}
	}

	// Launch 5 goroutines saving concurrently
	var wg sync.WaitGroup
	errors := make(chan error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := reg.Save(); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent Save error: %v", err)
	}

	// Verify file exists and is valid
	loaded, err := LoadRegistry(registryPath)
	if err != nil {
		t.Fatalf("LoadRegistry after concurrent saves failed: %v", err)
	}

	files := loaded.GetRegisteredFiles()
	if len(files) != 5 {
		t.Errorf("Expected 5 files in loaded registry, got %d", len(files))
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Register initial file
	if err := reg.RegisterFile("./test.txt", 0644, "user", "staff"); err != nil {
		t.Fatalf("RegisterFile failed: %v", err)
	}

	var wg sync.WaitGroup

	// 5 readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = reg.IsRegisteredFile("./test.txt")
				_, _ = reg.GetRegisteredFileMode("./test.txt")
				_ = reg.GetRegisteredFiles()
			}
		}()
	}

	// 5 writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				path := filepath.Join("./", "concurrent"+string(rune('0'+idx))+".txt")
				_ = reg.RegisterFile(path, 0644, "user", "staff")
			}
		}(i)
	}

	wg.Wait()

	// Verify registry is still consistent
	files := reg.GetRegisteredFiles()
	if len(files) < 1 {
		t.Error("Registry should have at least the initial file")
	}
}

// ============================================================================
// Test Category 2.6: Default Configuration Operations
// ============================================================================

func TestGetAndSetDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Test initial values
	if mode := reg.GetDefaultFileMode(); mode != 0640 {
		t.Errorf("Expected default mode 0640, got %o", mode)
	}
	if owner := reg.GetDefaultFileOwner(); owner != "root" {
		t.Errorf("Expected default owner 'root', got '%s'", owner)
	}
	if group := reg.GetDefaultFileGroup(); group != "wheel" {
		t.Errorf("Expected default group 'wheel', got '%s'", group)
	}

	// Set new values
	if err := reg.SetDefaultFileMode(0600); err != nil {
		t.Fatalf("SetDefaultFileMode failed: %v", err)
	}
	reg.SetDefaultFileOwner("newowner")
	reg.SetDefaultFileGroup("newgroup")

	// Verify new values
	if mode := reg.GetDefaultFileMode(); mode != 0600 {
		t.Errorf("Expected mode 0600, got %o", mode)
	}
	if owner := reg.GetDefaultFileOwner(); owner != "newowner" {
		t.Errorf("Expected owner 'newowner', got '%s'", owner)
	}
	if group := reg.GetDefaultFileGroup(); group != "newgroup" {
		t.Errorf("Expected group 'newgroup', got '%s'", group)
	}
}

func TestSetDefaultFileModeWithSpecialBits(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// NOTE: Registry design only stores permission bits (0-0777) using mode.Perm()
	// Special bits like sticky (01000) are masked out.
	// Setting mode 01000 (sticky bit) results in storing 0000 (permission bits only)
	err = reg.SetDefaultFileMode(01000)
	if err != nil {
		t.Fatalf("SetDefaultFileMode failed: %v", err)
	}

	// Verify mode was masked to permission bits only (01000.Perm() = 0000)
	if mode := reg.GetDefaultFileMode(); mode != 0000 {
		t.Errorf("Mode should be 0000 (permission bits only), got %o", mode)
	}
}

// ============================================================================
// Test Category: LastToggle Tracking
// ============================================================================

func TestLastToggle(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Initially empty
	name, toggleType := reg.GetLastToggle()
	if name != "" || toggleType != "" {
		t.Errorf("Expected empty last toggle, got name='%s' type='%s'", name, toggleType)
	}

	// Set last toggle
	reg.SetLastToggle("test.txt", "file")
	name, toggleType = reg.GetLastToggle()
	if name != "test.txt" {
		t.Errorf("Expected name 'test.txt', got '%s'", name)
	}
	if toggleType != "file" {
		t.Errorf("Expected type 'file', got '%s'", toggleType)
	}

	// Clear last toggle
	reg.ClearLastToggle()
	name, toggleType = reg.GetLastToggle()
	if name != "" || toggleType != "" {
		t.Errorf("Expected empty after clear, got name='%s' type='%s'", name, toggleType)
	}
}

func TestLastTogglePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "0640",
		GuardOwner: "root",
		GuardGroup: "wheel",
	}

	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// Set last toggle and save
	reg.SetLastToggle("mycoll", "collection")
	if err := reg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load and verify
	loaded, err := LoadRegistry(registryPath)
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	name, toggleType := loaded.GetLastToggle()
	if name != "mycoll" {
		t.Errorf("Expected name 'mycoll', got '%s'", name)
	}
	if toggleType != "collection" {
		t.Errorf("Expected type 'collection', got '%s'", toggleType)
	}
}
