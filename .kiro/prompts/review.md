---
description: "Review implemented code against specifications and requirements"
---

# Code Review Against Specifications

Perform a comprehensive code review of the implemented guard-tool functionality against the original specifications, steering documents, and shell test requirements.

## Review Process

1. **Load Context**
   - Review `.kiro/steering/` documents (product.md, tech.md, structure.md)
   - Analyze shell test specifications in `tests/` directory
   - Understand the 6-layer architecture requirements

2. **Architecture Compliance Review**
   - **CLI Layer**: Verify Cobra command structure matches `cmd/guard/commands/`
   - **Manager Layer**: Check business logic orchestration in `internal/manager/`
   - **Security Layer**: Validate path validation and symlink rejection
   - **Registry Layer**: Confirm YAML persistence integration
   - **Filesystem Layer**: Review platform-specific operations (macOS/Linux)
   - **TUI Layer**: Verify `-i` flag handling (if implemented)

3. **Specification Compliance**
   - **Auto-detection Logic**: Priority order (directory → file → collection → folder → registered file)
   - **Reserved Keywords**: Validation against 11 reserved words
   - **Exit Codes**: 0 for success/warnings, 1 for errors only
   - **Idempotent Operations**: Add existing files should skip, not error
   - **Privilege Handling**: Graceful warnings when not root, proper sudo usage in production

4. **Shell Test Compliance**
   - Compare implementation behavior against test expectations
   - Verify exact output formats, error messages, exit codes
   - Check edge cases (missing files, duplicate adds, invalid modes)
   - Validate collection operations match test specifications

5. **Code Quality Assessment**
   - **Go Best Practices**: Value semantics, error handling, file organization
   - **Function Length**: Maximum 50 lines per function
   - **Error Handling**: Proper wrapping with context, immediate checking
   - **Security**: Path validation, no directory traversal vulnerabilities
   - **Testing**: Coverage of critical protection logic

6. **Feature Completeness**
   - **Implemented Commands**: init, add, show, enable, disable, toggle, create, update
   - **Missing Commands**: remove, clear, destroy, config, cleanup, reset, uninstall, info, version
   - **Collection Support**: Creation, updates, protection operations
   - **File Operations**: Registration, protection, auto-detection

## Review Criteria

### ✅ **Must Have (Critical)**
- [ ] All shell tests pass with exact expected behavior
- [ ] 6-layer architecture properly implemented
- [ ] Auto-detection logic follows specified priority order
- [ ] Reserved keyword validation prevents conflicts
- [ ] Filesystem operations work with proper privilege checking
- [ ] Exit codes match specifications (0 for warnings, 1 for errors)

### ⚠️ **Should Have (Important)**
- [ ] Code follows Go best practices from tech.md
- [ ] Error messages match shell test expectations
- [ ] Idempotent operations behave correctly
- [ ] Security layer prevents path traversal attacks
- [ ] Cross-platform filesystem operations implemented

### 💡 **Nice to Have (Enhancement)**
- [ ] All 16 CLI commands implemented
- [ ] TUI mode with `-i` flag functional
- [ ] Comprehensive error handling and user feedback
- [ ] Performance optimizations for large file sets

## Deliverables

1. **Compliance Report**
   - Architecture adherence assessment
   - Shell test compliance status
   - Code quality evaluation
   - Security review findings

2. **Gap Analysis**
   - Missing functionality identification
   - Non-compliant implementations
   - Technical debt assessment

3. **Recommendations**
   - Priority fixes for critical issues
   - Suggestions for code improvements
   - Next steps for completion

## Focus Areas

- **Behavioral Accuracy**: Does the code do exactly what the tests expect?
- **Architectural Integrity**: Does the implementation follow the planned structure?
- **Security Robustness**: Are there any security vulnerabilities?
- **Code Maintainability**: Is the code well-organized and documented?
- **Specification Adherence**: Does every feature work as originally specified?

Execute this review systematically, documenting findings and providing specific recommendations for any issues discovered.
