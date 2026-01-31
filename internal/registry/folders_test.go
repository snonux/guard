package registry

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRegisterFolder(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	// Create a new registry
	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register a folder
	err = reg.RegisterFolder("@src/components", "src/components")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	// Verify the folder was registered
	folder := reg.GetFolder("@src/components")
	if folder == nil {
		t.Fatal("Expected folder to be registered, got nil")
	}
	if folder.Name != "@src/components" {
		t.Errorf("Expected folder name '@src/components', got '%s'", folder.Name)
	}
	if folder.Path != "src/components" {
		t.Errorf("Expected folder path 'src/components', got '%s'", folder.Path)
	}
	if folder.Guard != false {
		t.Error("Expected folder guard to be false by default")
	}
}

func TestRegisterFolderDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register a folder
	err = reg.RegisterFolder("@myfolder", "myfolder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	// Try to register the same folder again
	err = reg.RegisterFolder("@myfolder", "myfolder")
	if err == nil {
		t.Fatal("Expected error when registering duplicate folder, got nil")
	}
}

func TestUnregisterFolder(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register and then unregister a folder
	err = reg.RegisterFolder("@testfolder", "testfolder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	err = reg.UnregisterFolder("@testfolder", false)
	if err != nil {
		t.Fatalf("Failed to unregister folder: %v", err)
	}

	// Verify the folder is gone
	folder := reg.GetFolder("@testfolder")
	if folder != nil {
		t.Error("Expected folder to be unregistered, but it still exists")
	}
}

func TestUnregisterFolderNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Try to unregister a non-existent folder with ignoreMissing=false
	err = reg.UnregisterFolder("@nonexistent", false)
	if err == nil {
		t.Fatal("Expected error when unregistering non-existent folder, got nil")
	}

	// Try with ignoreMissing=true
	err = reg.UnregisterFolder("@nonexistent", true)
	if err != nil {
		t.Fatalf("Expected no error with ignoreMissing=true, got: %v", err)
	}
}

func TestGetFolder(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Get non-existent folder
	folder := reg.GetFolder("@nonexistent")
	if folder != nil {
		t.Error("Expected nil for non-existent folder")
	}

	// Register and get folder
	err = reg.RegisterFolder("@existing", "existing")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	folder = reg.GetFolder("@existing")
	if folder == nil {
		t.Fatal("Expected folder to exist")
	}
	if folder.Name != "@existing" {
		t.Errorf("Expected name '@existing', got '%s'", folder.Name)
	}
}

func TestGetFolderByPath(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Get by non-existent path
	folder := reg.GetFolderByPath("nonexistent")
	if folder != nil {
		t.Error("Expected nil for non-existent path")
	}

	// Register folder
	err = reg.RegisterFolder("@my/folder", "my/folder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	// Get by path
	folder = reg.GetFolderByPath("my/folder")
	if folder == nil {
		t.Fatal("Expected folder to exist by path")
	}
	if folder.Name != "@my/folder" {
		t.Errorf("Expected name '@my/folder', got '%s'", folder.Name)
	}
}

func TestSetFolderGuard(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Try to set guard on non-existent folder
	err = reg.SetFolderGuard("@nonexistent", true)
	if err == nil {
		t.Fatal("Expected error when setting guard on non-existent folder")
	}

	// Register folder
	err = reg.RegisterFolder("@testfolder", "testfolder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	// Set guard to true
	err = reg.SetFolderGuard("@testfolder", true)
	if err != nil {
		t.Fatalf("Failed to set folder guard: %v", err)
	}

	guard, err := reg.GetFolderGuard("@testfolder")
	if err != nil {
		t.Fatalf("Failed to get folder guard: %v", err)
	}
	if guard != true {
		t.Error("Expected folder guard to be true")
	}

	// Set guard back to false
	err = reg.SetFolderGuard("@testfolder", false)
	if err != nil {
		t.Fatalf("Failed to set folder guard: %v", err)
	}

	guard, err = reg.GetFolderGuard("@testfolder")
	if err != nil {
		t.Fatalf("Failed to get folder guard: %v", err)
	}
	if guard != false {
		t.Error("Expected folder guard to be false")
	}
}

func TestGetFolderGuard(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Try to get guard on non-existent folder
	_, err = reg.GetFolderGuard("@nonexistent")
	if err == nil {
		t.Fatal("Expected error when getting guard on non-existent folder")
	}

	// Register folder and check default guard
	err = reg.RegisterFolder("@testfolder", "testfolder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	guard, err := reg.GetFolderGuard("@testfolder")
	if err != nil {
		t.Fatalf("Failed to get folder guard: %v", err)
	}
	if guard != false {
		t.Error("Expected default folder guard to be false")
	}
}

func TestListFolders(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// List empty folders
	folders := reg.ListFolders()
	if len(folders) != 0 {
		t.Errorf("Expected 0 folders, got %d", len(folders))
	}

	// Register some folders
	err = reg.RegisterFolder("@folder1", "folder1")
	if err != nil {
		t.Fatalf("Failed to register folder1: %v", err)
	}
	err = reg.RegisterFolder("@folder2", "folder2")
	if err != nil {
		t.Fatalf("Failed to register folder2: %v", err)
	}

	folders = reg.ListFolders()
	if len(folders) != 2 {
		t.Errorf("Expected 2 folders, got %d", len(folders))
	}
}

func TestIsRegisteredFolder(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Check non-existent folder
	if reg.IsRegisteredFolder("@nonexistent") {
		t.Error("Expected false for non-existent folder")
	}

	// Register folder and check
	err = reg.RegisterFolder("@myfolder", "myfolder")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	if !reg.IsRegisteredFolder("@myfolder") {
		t.Error("Expected true for registered folder")
	}
}

func TestIsRegisteredFolderByPath(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Check non-existent path
	if reg.IsRegisteredFolderByPath("nonexistent") {
		t.Error("Expected false for non-existent path")
	}

	// Register folder and check by path
	err = reg.RegisterFolder("@my/path", "my/path")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}

	if !reg.IsRegisteredFolderByPath("my/path") {
		t.Error("Expected true for registered folder path")
	}
}

func TestCountFolders(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Count empty
	if count := reg.CountFolders(); count != 0 {
		t.Errorf("Expected 0 folders, got %d", count)
	}

	// Register folders and count
	if err := reg.RegisterFolder("@f1", "f1"); err != nil {
		t.Fatalf("Failed to register folder f1: %v", err)
	}
	if err := reg.RegisterFolder("@f2", "f2"); err != nil {
		t.Fatalf("Failed to register folder f2: %v", err)
	}
	if err := reg.RegisterFolder("@f3", "f3"); err != nil {
		t.Fatalf("Failed to register folder f3: %v", err)
	}

	if count := reg.CountFolders(); count != 3 {
		t.Errorf("Expected 3 folders, got %d", count)
	}
}

func TestGetRegisteredFolders(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Get empty list
	names := reg.GetRegisteredFolders()
	if len(names) != 0 {
		t.Errorf("Expected 0 folder names, got %d", len(names))
	}

	// Register folders
	if err := reg.RegisterFolder("@alpha", "alpha"); err != nil {
		t.Fatalf("Failed to register folder alpha: %v", err)
	}
	if err := reg.RegisterFolder("@beta", "beta"); err != nil {
		t.Fatalf("Failed to register folder beta: %v", err)
	}

	names = reg.GetRegisteredFolders()
	if len(names) != 2 {
		t.Errorf("Expected 2 folder names, got %d", len(names))
	}

	// Check that both names are present
	found := make(map[string]bool)
	for _, name := range names {
		found[name] = true
	}
	if !found["@alpha"] || !found["@beta"] {
		t.Error("Expected both @alpha and @beta in folder names")
	}
}

func TestFolderYAMLSerialization(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register a folder with guard enabled
	err = reg.RegisterFolder("@src/components", "src/components")
	if err != nil {
		t.Fatalf("Failed to register folder: %v", err)
	}
	err = reg.SetFolderGuard("@src/components", true)
	if err != nil {
		t.Fatalf("Failed to set folder guard: %v", err)
	}

	// Save to disk
	err = reg.Save()
	if err != nil {
		t.Fatalf("Failed to save registry: %v", err)
	}

	// Read the file and verify YAML structure
	data, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	var registryData RegistryData
	err = yaml.Unmarshal(data, &registryData)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(registryData.Folders) != 1 {
		t.Fatalf("Expected 1 folder in YAML, got %d", len(registryData.Folders))
	}

	folder := registryData.Folders[0]
	if folder.Name != "@src/components" {
		t.Errorf("Expected folder name '@src/components', got '%s'", folder.Name)
	}
	if folder.Path != "src/components" {
		t.Errorf("Expected folder path 'src/components', got '%s'", folder.Path)
	}
	if folder.Guard != true {
		t.Error("Expected folder guard to be true in YAML")
	}
}

func TestFolderYAMLDeserialization(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	// Write a YAML file with folders section
	yamlContent := `config:
  guard_mode: "0000"
  guard_owner: testuser
  guard_group: testgroup
files: []
collections: []
folders:
  - name: "@internal/registry"
    path: "internal/registry"
    guard: true
  - name: "@cmd/guard"
    path: "cmd/guard"
    guard: false
`
	err := os.WriteFile(registryPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	// Load the registry
	reg, err := LoadRegistry(registryPath)
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	// Verify folders were loaded
	if reg.CountFolders() != 2 {
		t.Errorf("Expected 2 folders, got %d", reg.CountFolders())
	}

	// Check first folder
	folder1 := reg.GetFolder("@internal/registry")
	if folder1 == nil {
		t.Fatal("Expected @internal/registry folder to exist")
	}
	if folder1.Path != "internal/registry" {
		t.Errorf("Expected path 'internal/registry', got '%s'", folder1.Path)
	}
	if folder1.Guard != true {
		t.Error("Expected folder1 guard to be true")
	}

	// Check second folder
	folder2 := reg.GetFolder("@cmd/guard")
	if folder2 == nil {
		t.Fatal("Expected @cmd/guard folder to exist")
	}
	if folder2.Path != "cmd/guard" {
		t.Errorf("Expected path 'cmd/guard', got '%s'", folder2.Path)
	}
	if folder2.Guard != false {
		t.Error("Expected folder2 guard to be false")
	}
}

func TestFolderLoadAfterSave(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	// Create registry with folders
	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register folders
	if err := reg.RegisterFolder("@folder1", "folder1"); err != nil {
		t.Fatalf("Failed to register folder1: %v", err)
	}
	if err := reg.SetFolderGuard("@folder1", true); err != nil {
		t.Fatalf("Failed to set folder1 guard: %v", err)
	}
	if err := reg.RegisterFolder("@folder2", "folder2"); err != nil {
		t.Fatalf("Failed to register folder2: %v", err)
	}

	// Save
	err = reg.Save()
	if err != nil {
		t.Fatalf("Failed to save registry: %v", err)
	}

	// Load into new registry
	reg2, err := LoadRegistry(registryPath)
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	// Verify folders persist
	if reg2.CountFolders() != 2 {
		t.Errorf("Expected 2 folders after reload, got %d", reg2.CountFolders())
	}

	guard1, _ := reg2.GetFolderGuard("@folder1")
	if guard1 != true {
		t.Error("Expected folder1 guard to persist as true")
	}

	guard2, _ := reg2.GetFolderGuard("@folder2")
	if guard2 != false {
		t.Error("Expected folder2 guard to persist as false")
	}
}

func TestGetFolderReturnsCopy(t *testing.T) {
	tmpDir := t.TempDir()
	registryPath := filepath.Join(tmpDir, ".guardfile")

	defaults := &RegistryDefaults{
		GuardMode:  "000",
		GuardOwner: "testuser",
		GuardGroup: "testgroup",
	}
	reg, err := NewRegistry(registryPath, defaults, false)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Register folder
	if err := reg.RegisterFolder("@test", "test"); err != nil {
		t.Fatalf("Failed to register folder test: %v", err)
	}

	// Get folder and modify it
	folder := reg.GetFolder("@test")
	folder.Guard = true
	folder.Path = "modified"

	// Get again and verify original is unchanged
	folder2 := reg.GetFolder("@test")
	if folder2.Guard != false {
		t.Error("Expected original folder guard to remain false (copy returned)")
	}
	if folder2.Path != "test" {
		t.Error("Expected original folder path to remain 'test' (copy returned)")
	}
}
