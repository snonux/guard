# Project Structure

## Directory Layout
```
guard-tool/
├── cmd/
│   └── guard/
│       ├── main.go        # Entry point
│       └── commands/      # CLI commands
│           ├── init.go
│           ├── add.go
│           ├── remove.go
│           ├── enable.go
│           ├── disable.go
│           ├── toggle.go
│           ├── create.go
│           ├── update.go
│           ├── clear.go
│           ├── destroy.go
│           ├── show.go
│           ├── info.go
│           ├── config.go
│           ├── cleanup.go
│           ├── reset.go
│           ├── uninstall.go
│           └── version.go
├── internal/              # Private application code
│   ├── manager/           # Business logic orchestration
│   ├── security/          # Path validation and symlink rejection
│   ├── registry/          # YAML persistence and state management
│   ├── filesystem/        # OS operations (chmod, chown, immutable flags)
│   └── tui/               # Bubble Tea interface
├── tests/                 # Shell integration tests
├── docs/                  # Documentation
├── .guardfile             # YAML state file (created by init)
├── justfile               # Build and task automation
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## File Naming Conventions
- **Go files**: snake_case.go (following Go conventions)
- **Commands**: Single word commands in commands/ directory
- **Packages**: Short, lowercase names without underscores
- **Constants**: Standard Go conventions with iota for related constants
- **State file**: .guardfile (dot prefix for hidden file)

## Module Organization
- **cmd/guard/**: Entry point and command structure
- **cmd/guard/commands/**: Individual CLI commands using Cobra patterns
- **internal/manager/**: Business logic orchestration and workflow coordination
- **internal/security/**: Path validation and symlink rejection
- **internal/registry/**: YAML parsing, state persistence, validation
- **internal/filesystem/**: Unix syscalls, filesystem operations, sudo handling
- **internal/tui/**: Bubble Tea models, views, and update functions

## Configuration Files
- **.guardfile**: YAML state file tracking protected files and collections
- **go.mod/go.sum**: Go module dependencies
- **justfile**: Build tasks, testing, cross-compilation scripts
- **.gitignore**: Exclude build artifacts and temporary files

## Documentation Structure
- **README.md**: Installation, quick start, basic usage
- **docs/**: Detailed documentation, examples, troubleshooting
- **Godoc comments**: Inline documentation for all exported functions

## Build Artifacts
- **GOPATH/bin/**: Compiled binaries installed via `just install`
- **coverage.out**: Test coverage reports

## Environment-Specific Files
- **Development**: Local .guardfile for testing
