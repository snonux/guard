# Feature: Refactor Filesystem Component

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the filesystem component from interface-based design with build tags to a single concrete struct with runtime OS detection. Add comprehensive filesystem operations needed by the manager layer and improve platform-specific immutable flag handling.

## User Story

As a developer maintaining the guard-tool codebase
I want a simplified filesystem component with all required operations
So that the manager layer can perform file operations efficiently without missing functionality

## Problem Statement

The current filesystem implementation has several issues:
1. Interface-based design with separate platform files creates unnecessary complexity
2. Missing critical methods needed by the manager layer (FileExists, GetFileInfo, etc.)
3. Immutable flag operations don't preserve existing flags on Darwin
4. Linux implementation uses unsafe directly instead of proper unix package methods
5. HasRootPrivileges is a package function instead of struct method

## Solution Statement

Replace the interface-based design with a single concrete FileSystem struct that uses runtime.GOOS checks for platform-specific operations. Add all missing methods required by the manager layer and improve platform-specific implementations.

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: internal/filesystem, internal/manager
**Dependencies**: golang.org/x/sys/unix

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/filesystem/filesystem.go` (lines 1-54) - Why: Current interface design to be replaced
- `internal/filesystem/filesystem_darwin.go` (lines 1-76) - Why: Darwin-specific implementation patterns
- `internal/filesystem/filesystem_linux.go` (lines 1-107) - Why: Linux-specific implementation patterns  
- `internal/manager/files.go` (lines 60-120) - Why: Current filesystem usage patterns in manager
- `internal/manager/manager.go` (lines 1-50) - Why: Manager struct and filesystem integration
- `go.mod` - Why: Current dependencies and Go version

### New Files to Create

- `internal/filesystem/filesystem.go` - Single consolidated filesystem implementation

### Files to Remove

- `internal/filesystem/filesystem_darwin.go` - Replaced by runtime checks
- `internal/filesystem/filesystem_linux.go` - Replaced by runtime checks

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [golang.org/x/sys/unix Documentation](https://pkg.go.dev/golang.org/x/sys/unix)
  - Specific section: IoctlGetInt and IoctlSetInt functions
  - Why: Required for proper Linux immutable flag handling
- [Go runtime Package](https://pkg.go.dev/runtime)
  - Specific section: GOOS constant
  - Why: Runtime OS detection for platform-specific code

### Patterns to Follow

**Error Handling Pattern:**
```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}
```

**Privilege Check Pattern:**
```go
if !HasRootPrivileges() {
    fmt.Printf("Warning: Skipping operation (requires root privileges): %s\n", path)
    return nil
}
```

**Manager Integration Pattern:**
```go
// From manager.go line 32
fs: filesystem.NewFilesystem(),
```

**Naming Conventions:**
- Struct methods use receiver `(fs *FileSystem)`
- Public methods are PascalCase
- Private helper functions are camelCase
- Error messages start with lowercase and use %w for wrapping

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Create the new consolidated filesystem implementation with runtime OS detection and all required methods.

**Tasks:**
- Remove interface design and create concrete FileSystem struct
- Implement runtime.GOOS-based platform detection
- Add all missing methods required by manager layer

### Phase 2: Platform-Specific Implementation

Implement improved platform-specific code for immutable flags and ownership operations.

**Tasks:**
- Implement Darwin immutable flag operations with proper flag preservation
- Implement Linux immutable flag operations using unix.IoctlGetInt/IoctlSetInt
- Consolidate ownership operations with proper error handling

### Phase 3: Integration

Update manager layer to use new filesystem implementation and remove old files.

**Tasks:**
- Update manager imports and usage
- Remove old platform-specific files
- Verify all existing functionality works

### Phase 4: Testing & Validation

Ensure all operations work correctly on both platforms.

**Tasks:**
- Test immutable flag operations preserve existing flags
- Test all new methods work correctly
- Validate manager integration

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE internal/filesystem/filesystem.go

- **IMPLEMENT**: Complete FileSystem struct with all required methods
- **PATTERN**: Error handling from existing files (filesystem_*.go)
- **IMPORTS**: os, fmt, runtime, os/user, strconv, golang.org/x/sys/unix
- **GOTCHA**: Must handle both Darwin and Linux platforms in single file
- **VALIDATE**: `go build ./internal/filesystem`

```go
package filesystem

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

// FileSystem provides filesystem operations with platform-specific implementations
type FileSystem struct{}

// NewFileSystem creates a new FileSystem instance
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// HasRootPrivileges checks if running as root
func (fs *FileSystem) HasRootPrivileges() bool {
	return os.Geteuid() == 0
}

// FileExists checks if path exists and is a regular file (not directory)
func (fs *FileSystem) FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// GetFileInfo returns file mode, owner, and group together
func (fs *FileSystem) GetFileInfo(path string) (mode os.FileMode, owner, group string, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to stat file: %w", err)
	}
	
	mode = info.Mode().Perm() // Only permission bits
	
	// Get owner and group from system info
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if u, err := user.LookupId(strconv.Itoa(int(stat.Uid))); err == nil {
			owner = u.Username
		} else {
			owner = strconv.Itoa(int(stat.Uid))
		}
		
		if g, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid))); err == nil {
			group = g.Name
		} else {
			group = strconv.Itoa(int(stat.Gid))
		}
	}
	
	return mode, owner, group, nil
}

// CheckFilesExist partitions paths into existing and missing slices
func (fs *FileSystem) CheckFilesExist(paths []string) (existing, missing []string) {
	for _, path := range paths {
		if fs.FileExists(path) {
			existing = append(existing, path)
		} else {
			missing = append(missing, path)
		}
	}
	return existing, missing
}

// ApplyPermissions applies mode, owner, and group in sequence
func (fs *FileSystem) ApplyPermissions(path string, mode os.FileMode, owner, group string) error {
	if err := fs.Chmod(path, mode); err != nil {
		return fmt.Errorf("failed to set mode: %w", err)
	}
	
	if err := fs.Chown(path, owner); err != nil {
		return fmt.Errorf("failed to set owner: %w", err)
	}
	
	if err := fs.Chgrp(path, group); err != nil {
		return fmt.Errorf("failed to set group: %w", err)
	}
	
	return nil
}

// RestorePermissions restores mode, owner, and group in sequence
func (fs *FileSystem) RestorePermissions(path string, mode os.FileMode, owner, group string) error {
	// Same as ApplyPermissions but semantically different
	return fs.ApplyPermissions(path, mode, owner, group)
}

// Chmod changes file permissions
func (fs *FileSystem) Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// Chown changes file owner only
func (fs *FileSystem) Chown(path string, owner string) error {
	if !fs.HasRootPrivileges() {
		fmt.Printf("Warning: Skipping ownership change (requires root privileges): %s\n", path)
		return nil
	}

	// Look up user ID
	u, err := user.Lookup(owner)
	if err != nil {
		return fmt.Errorf("failed to lookup user %s: %w", owner, err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := os.Chown(path, uid, -1); err != nil {
		return fmt.Errorf("failed to change owner: %w", err)
	}
	return nil
}

// Chgrp changes file group only
func (fs *FileSystem) Chgrp(path string, group string) error {
	if !fs.HasRootPrivileges() {
		fmt.Printf("Warning: Skipping group change (requires root privileges): %s\n", path)
		return nil
	}

	g, err := user.LookupGroup(group)
	if err != nil {
		return fmt.Errorf("failed to lookup group %s: %w", group, err)
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	if err := os.Chown(path, -1, gid); err != nil {
		return fmt.Errorf("failed to change group: %w", err)
	}
	return nil
}

// DirEntry represents a directory entry
type DirEntry struct {
	Name     string
	Path     string
	IsDir    bool
	IsLink   bool
	FileInfo os.FileInfo
}

// ReadDir returns directory entries with detailed information, sorted with directories first
func (fs *FileSystem) ReadDir(dirPath string) ([]DirEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var result []DirEntry
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		
		// Use Lstat to detect symlinks
		lstatInfo, err := os.Lstat(fullPath)
		if err != nil {
			continue // Skip entries we can't lstat
		}

		isLink := lstatInfo.Mode()&os.ModeSymlink != 0
		var isDir bool
		var fileInfo os.FileInfo

		if isLink {
			// For symlinks, use Stat to determine if target is directory
			statInfo, err := os.Stat(fullPath)
			if err != nil {
				// Broken symlink, use lstat info
				isDir = false
				fileInfo = lstatInfo
			} else {
				isDir = statInfo.IsDir()
				fileInfo = statInfo
			}
		} else {
			isDir = lstatInfo.IsDir()
			fileInfo = lstatInfo
		}

		dirEntry := DirEntry{
			Name:     entry.Name(),
			Path:     fullPath,
			IsDir:    isDir,
			IsLink:   isLink,
			FileInfo: fileInfo,
		}
		result = append(result, dirEntry)
	}

	// Sort with directories first, then alphabetically by name
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir // Directories first
		}
		return result[i].Name < result[j].Name // Alphabetical within groups
	})

	return result, nil
}

// Lstat returns file info without following symlinks
func (fs *FileSystem) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

// IsSymlink checks if path is a symbolic link
func (fs *FileSystem) IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, fmt.Errorf("failed to lstat: %w", err)
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

// IsDir checks if path is a directory
func (fs *FileSystem) IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat: %w", err)
	}
	return info.IsDir(), nil
}

// CollectImmediateFiles returns regular files in directory (non-recursive)
func (fs *FileSystem) CollectImmediateFiles(dirPath string) ([]string, error) {
	entries, err := fs.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir && !entry.IsLink && entry.FileInfo.Mode().IsRegular() {
			files = append(files, entry.Path)
		}
	}

	return files, nil
}

// CollectFilesRecursive returns all regular files in directory tree
func (fs *FileSystem) CollectFilesRecursive(dirPath string) ([]string, error) {
	var files []string
	
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// Skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		
		// Add regular files
		if d.Type().IsRegular() {
			files = append(files, path)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	return files, nil
}

// SetImmutable sets the immutable flag on a file
func (fs *FileSystem) SetImmutable(path string) error {
	if !fs.HasRootPrivileges() {
		fmt.Printf("Warning: Skipping immutable flag change (requires root privileges): %s\n", path)
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return fs.setImmutableDarwin(path, true)
	case "linux":
		return fs.setImmutableLinux(path, true)
	default:
		return fmt.Errorf("immutable flags not supported on %s", runtime.GOOS)
	}
}

// ClearImmutable removes the immutable flag from a file
func (fs *FileSystem) ClearImmutable(path string) error {
	if !fs.HasRootPrivileges() {
		fmt.Printf("Warning: Skipping immutable flag change (requires root privileges): %s\n", path)
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return fs.setImmutableDarwin(path, false)
	case "linux":
		return fs.setImmutableLinux(path, false)
	default:
		return fmt.Errorf("immutable flags not supported on %s", runtime.GOOS)
	}
}

// IsImmutable checks if file has immutable flag set
func (fs *FileSystem) IsImmutable(path string) (bool, error) {
	switch runtime.GOOS {
	case "darwin":
		return fs.isImmutableDarwin(path)
	case "linux":
		return fs.isImmutableLinux(path)
	default:
		return false, fmt.Errorf("immutable flags not supported on %s", runtime.GOOS)
	}
}

// Darwin-specific immutable flag operations
func (fs *FileSystem) setImmutableDarwin(path string, immutable bool) error {
	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	flags := stat.Flags
	if immutable {
		flags |= unix.SF_IMMUTABLE // OR with existing flags
	} else {
		flags &^= unix.SF_IMMUTABLE // AND NOT to clear flag
	}

	if err := unix.Chflags(path, int(flags)); err != nil {
		return fmt.Errorf("failed to set immutable flag: %w", err)
	}
	return nil
}

func (fs *FileSystem) isImmutableDarwin(path string) (bool, error) {
	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err != nil {
		return false, fmt.Errorf("failed to get file stats: %w", err)
	}
	return (stat.Flags & unix.SF_IMMUTABLE) != 0, nil
}

// Linux-specific immutable flag operations
const (
	FS_IOC_GETFLAGS = 0x80086601
	FS_IOC_SETFLAGS = 0x40086602
	FS_IMMUTABLE_FL = 0x00000010
)

func (fs *FileSystem) setImmutableLinux(path string, immutable bool) error {
	fd, err := unix.Open(path, unix.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer unix.Close(fd)

	// Get current flags using proper unix package method
	flags, err := unix.IoctlGetInt(fd, FS_IOC_GETFLAGS)
	if err != nil {
		return fmt.Errorf("failed to get file flags: %w", err)
	}

	// Set or unset immutable flag
	if immutable {
		flags |= FS_IMMUTABLE_FL
	} else {
		flags &^= FS_IMMUTABLE_FL
	}

	// Set new flags using proper unix package method
	if err := unix.IoctlSetInt(fd, FS_IOC_SETFLAGS, flags); err != nil {
		return fmt.Errorf("failed to set file flags: %w", err)
	}

	return nil
}

func (fs *FileSystem) isImmutableLinux(path string) (bool, error) {
	fd, err := unix.Open(path, unix.O_RDONLY, 0)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer unix.Close(fd)

	// Get current flags using proper unix package method
	flags, err := unix.IoctlGetInt(fd, FS_IOC_GETFLAGS)
	if err != nil {
		return false, fmt.Errorf("failed to get file flags: %w", err)
	}

	return (flags & FS_IMMUTABLE_FL) != 0, nil
}

// Legacy methods for backward compatibility
func (fs *FileSystem) SetPermissions(path string, mode os.FileMode) error {
	return fs.Chmod(path, mode)
}

func (fs *FileSystem) SetOwnership(path string, owner, group string) error {
	if err := fs.Chown(path, owner); err != nil {
		return err
	}
	return fs.Chgrp(path, group)
}

func (fs *FileSystem) GetPermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}
```

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Update import and struct field type
- **PATTERN**: Existing manager initialization pattern (line 32)
- **IMPORTS**: Change filesystem import usage
- **GOTCHA**: Must update NewFilesystem() call to NewFileSystem()
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/files.go

- **IMPLEMENT**: Update filesystem method calls to use new methods
- **PATTERN**: Existing error handling in enableFileProtection/disableFileProtection
- **IMPORTS**: No import changes needed
- **GOTCHA**: SetImmutable(path, true) becomes SetImmutable(path), SetImmutable(path, false) becomes ClearImmutable(path)
- **VALIDATE**: `go build ./internal/manager`

### REMOVE internal/filesystem/filesystem_darwin.go

- **IMPLEMENT**: Delete the file
- **PATTERN**: N/A
- **IMPORTS**: N/A
- **GOTCHA**: Ensure no other files import this directly
- **VALIDATE**: `find . -name "*.go" -exec grep -l "filesystem_darwin" {} \;` should return empty

### REMOVE internal/filesystem/filesystem_linux.go

- **IMPLEMENT**: Delete the file
- **PATTERN**: N/A
- **IMPORTS**: N/A
- **GOTCHA**: Ensure no other files import this directly
- **VALIDATE**: `find . -name "*.go" -exec grep -l "filesystem_linux" {} \;` should return empty

### UPDATE go.mod dependencies

- **IMPLEMENT**: Ensure golang.org/x/sys dependency is present
- **PATTERN**: Existing require block in go.mod
- **IMPORTS**: N/A
- **GOTCHA**: Version should be compatible with existing usage
- **VALIDATE**: `go mod tidy && go mod verify`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests for each new method following Go testing conventions:
- Test FileExists with files and directories
- Test GetFileInfo return values
- Test CheckFilesExist partitioning logic
- Test ApplyPermissions/RestorePermissions sequences
- Test ReadDir with various directory contents
- Test CollectImmediateFiles vs CollectFilesRecursive
- Test SetImmutable/ClearImmutable/IsImmutable operations

### Integration Tests

Test filesystem operations with manager layer:
- Test manager can enable/disable file protection
- Test immutable flag operations work end-to-end
- Test permission restoration works correctly

### Edge Cases

Test specific edge cases for this feature:
- Test immutable flag preservation on Darwin
- Test proper error handling when not root
- Test symlink handling in directory operations
- Test missing file handling in CheckFilesExist

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./internal/filesystem/
go vet ./internal/filesystem/
go build ./internal/filesystem/
```

### Level 2: Unit Tests

```bash
go test ./internal/filesystem/ -v
go test ./internal/manager/ -v
```

### Level 3: Integration Tests

```bash
go build ./cmd/guard/
./tests/run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test basic functionality
./guard init 644 $(whoami) $(id -gn)
echo "test" > testfile.txt
./guard add testfile.txt
./guard enable testfile.txt
ls -la testfile.txt
./guard disable testfile.txt
rm testfile.txt .guardfile
```

### Level 5: Platform-Specific Validation

```bash
# Test immutable flags (requires root)
sudo ./guard init 644 root root
sudo echo "test" > testfile.txt
sudo ./guard add testfile.txt
sudo ./guard enable testfile.txt
# On Darwin: ls -lO testfile.txt should show "schg"
# On Linux: lsattr testfile.txt should show "i"
sudo ./guard disable testfile.txt
sudo rm testfile.txt .guardfile
```

---

## ACCEPTANCE CRITERIA

- [ ] Single filesystem.go file replaces interface design
- [ ] All new methods (FileExists, GetFileInfo, etc.) implemented
- [ ] Runtime GOOS detection works for Darwin and Linux
- [ ] Darwin SetImmutable preserves existing flags using OR/AND NOT
- [ ] Linux implementation uses unix.IoctlGetInt/IoctlSetInt
- [ ] HasRootPrivileges is a method on FileSystem struct
- [ ] Manager layer integration works without changes to external API
- [ ] All existing tests pass
- [ ] No build tag files remain
- [ ] Backward compatibility maintained for existing manager usage

---

## COMPLETION CHECKLIST

- [ ] New filesystem.go created with all required methods
- [ ] Platform-specific implementations use runtime.GOOS
- [ ] Darwin immutable operations preserve existing flags
- [ ] Linux immutable operations use proper unix package methods
- [ ] Manager updated to use new FileSystem struct
- [ ] Old platform-specific files removed
- [ ] All validation commands pass
- [ ] Integration tests confirm functionality
- [ ] No regressions in existing behavior

---

## NOTES

**Key Design Decisions:**
- Single struct approach simplifies maintenance and testing
- Runtime GOOS checks provide platform flexibility without build complexity
- Separate SetImmutable/ClearImmutable methods provide clearer API
- DirEntry struct provides rich directory information
- Legacy method compatibility ensures smooth transition

**Implementation Risks:**
- Platform-specific code must be tested on both Darwin and Linux
- Immutable flag operations require root privileges
- Directory traversal operations must handle symlinks correctly
- Error handling must be consistent across all new methods

**Confidence Score**: 8/10 - Well-defined requirements with clear patterns from existing code, but platform-specific testing required for full validation.
