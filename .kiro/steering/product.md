# Product Overview

## Product Purpose
Guard-tool protects files from unwanted modifications by AI coding agents like Cursor, Claude, Kiro, Codex, and GitHub Copilot. It provides a convenient way to quickly toggle file protection at the press of a button, making files truly read-only through filesystem-level permissions and immutable flags.

## Target Users
Developers who use AI coding assistants and want to prevent those tools from modifying files outside the scope of their current work. These users need:
- Quick protection for critical files (configs, docs, completed features)
- Easy toggle between protected/unprotected states
- Confidence that AI tools won't accidentally modify important files
- Minimal workflow disruption when switching between different coding tasks

## Key Features
- **File Protection**: Modify rwx permissions, owner, group, and immutable flags (requires sudo)
- **Quick Toggle**: One-command protection/unprotection of files and directories
- **Interactive TUI**: Bubble Tea-powered terminal interface for file selection
- **Collections**: Named groups of files that can be toggled together
- **Folders**: Dynamic collections that scan files from directory paths
- **State Persistence**: YAML-based .guardfile for tracking protected files
- **Security Validation**: Path validation and symlink rejection to prevent exploits
- **Cross-Platform**: Support for Linux, macOS, and BSD systems
- **CLI Interface**: Cobra-powered command-line interface with intuitive commands

## Business Objectives
- Provide developers with confidence when using AI coding tools
- Reduce accidental file modifications and related debugging time
- Enable safer AI-assisted development workflows
- Build a reliable, fast, and user-friendly protection system

## User Journey
1. **Setup**: Run `guard init` with mode, owner, and group parameters to create .guardfile
2. **Register**: Use `guard add` to register files and capture their current permissions
3. **Protection**: Enable protection on registered files using CLI or TUI
4. **Development**: Work with AI tools knowing protected files are safe
5. **Toggle**: Quickly unprotect files when they need modification
6. **Collection Management**: Manage protection states for named groups and folders

## Success Criteria
- Zero accidental modifications to protected files by AI tools
- Sub-second response time for protection/unprotection operations
- Intuitive workflow that doesn't slow down development
- Reliable state persistence across system restarts
- Positive developer feedback on workflow integration
