# Technical Architecture

## Technology Stack
- **Primary Language**: Go (for performance, cross-platform support, and excellent CLI tooling)
- **CLI Framework**: Cobra (industry standard for Go command-line applications)
- **Terminal UI**: Bubble Tea (reactive model for interactive terminal interfaces)
- **Build System**: Just (modern command runner for task automation)
- **Configuration**: YAML (.guardfile for state persistence)
- **Target Platforms**: Unix-like systems (Linux, macOS, BSD)

## Architecture Overview
- **CLI Layer**: Cobra-based command parsing and argument handling
- **Manager Layer**: Business logic orchestration and workflow coordination
- **Security Layer**: Path validation and symlink rejection
- **Registry Layer**: YAML persistence and state management for .guardfile
- **Filesystem Layer**: OS operations (chmod, chown, immutable flags)
- **TUI Layer**: Bubble Tea interactive interface for file selection

## Development Environment
- **Go Version**: Latest stable Go release
- **Package Manager**: Go modules (go.mod/go.sum)
- **Build Tool**: Just command runner with justfile
- **Dependencies**: Cobra CLI, Bubble Tea TUI, YAML parser
- **Development Tools**: gofmt, golangci-lint, go vet for code quality

## Code Standards
- **Formatting**: Standard Go formatting with gofmt
- **Naming**: Go naming conventions (PascalCase for exported, camelCase for unexported)
- **Documentation**: Godoc comments for all exported functions and types
- **Error Handling**: Explicit error handling with wrapped errors for context
- **Package Structure**: Clear separation between CLI, TUI, core logic, and platform layers
- **Function Length**: Maximum 50 lines per function
- **File Organization**: Constants and types at top, constructors after type definitions, public functions before private
- **Value Semantics**: Use value semantics unless pointers are necessary
- **Receiver Types**: Consistent receiver types (all pointer or all value, never mixed)
- **Context Usage**: context.Context as first parameter for blocking or IO operations
- **Dependency Injection**: Prefer dependency injection over package-level variables
- **Resource Management**: Use defer to close resources immediately after opening
- **Constants**: Use iota for related constants
- **Error Recovery**: Avoid panic except for truly unrecoverable errors
- **Code Structure**: Avoid deep nesting, prefer composition over inheritance
- **CI Pipeline**: The CI pipeline MUST pass with exit code 0 - any CI failures are unacceptable and block development

## Error Handling Standards
- **Return Values**: Errors returned as last return value
- **Error Checking**: Errors checked immediately after calls
- **Error Wrapping**: Wrap errors with context using fmt.Errorf with %w
- **Error Types**: Use errors.Is and errors.As for checking error types
- **CLI Exit Codes**: Exit with code 1 on error, 0 on success
- **Library Code**: Return errors, never call os.Exit or log.Fatal

## Testing Strategy
- **Unit Tests**: Go's built-in testing framework for core logic
- **Table-Driven Tests**: Unit tests use table-driven patterns
- **Integration Tests**: Shell integration tests in tests/ directory with run-all-tests script
- **TUI Tests**: Interactive tests requiring tmux for terminal simulation
- **Permission Tests**: Validation of filesystem permission and immutable flag operations
- **Sudo-Free Testing**: Tests simulate permission failures using directory permissions instead of requiring sudo
- **Coverage Target**: 60%+ test coverage for critical protection logic

Note: Runtime requires sudo for full protection (ownership changes, immutable flags), but the test suite runs without sudo by using directory permission manipulation to simulate access failures. This enables CI/CD pipelines and developer testing without elevated privileges.

**Shell Integration Tests**: The tests/ directory contains shell-based integration tests that serve as both behavioral tests and executable specifications. These tests define the exact expected output formats, error messages, warning messages, and exit codes for all CLI scenarios. When implementing CLI commands, reference these shell tests to understand precise expected behavior, including edge cases like adding duplicate files or enabling protection on non-existent files. Tests cover file operations, collection operations, folder operations, configuration commands, and maintenance commands.

## Deployment Process
- **Build**: Cross-compilation for multiple Unix platforms using Go build
- **Installation**: Clone repository, run `just build`, then `just install` to place binary in GOPATH/bin
- **Distribution**: Source code distribution via Git repository
- **Versioning**: Semantic versioning with Git tags

## Performance Requirements
- **Response Time**: Sub-second file protection/unprotection operations
- **Memory Usage**: Minimal memory footprint for CLI operations
- **Startup Time**: Near-instantaneous command execution
- **File Handling**: Efficient batch operations for large file collections

## Security Considerations
- **Privilege Escalation**: Secure sudo usage for filesystem operations
- **File Permissions**: Safe handling of rwx permissions and ownership changes
- **Tampering Detection**: Path validation on registry load to detect tampering
- **Input Validation**: Path validation and symlink rejection to prevent exploits
