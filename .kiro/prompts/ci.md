---
description: "Run CI pipeline and fix any issues found"
---

# Fix CI Issues

Run the complete CI pipeline and systematically fix any issues that are discovered.

## Process

1. **Run CI Pipeline**
   ```bash
   just ci
   ```

2. **Analyze Results**
   - Review all linting errors, test failures, and build issues
   - Categorize issues by type (formatting, linting, tests, build)
   - Prioritize fixes by impact and complexity

3. **Fix Issues Systematically**
   - **Formatting Issues**: Run `go fmt ./...` and `goimports` fixes
   - **Linting Issues**: Address golangci-lint, gocyclo, gocognit warnings
   - **Security Issues**: Fix gosec and semgrep findings
   - **Test Failures**: Debug and fix failing shell integration tests
   - **Build Issues**: Resolve compilation errors and dependency problems

4. **Validate Fixes**
   - Run `just ci` after each batch of fixes
   - Ensure no regressions are introduced
   - Verify all tests pass

5. **Final Verification**
   - Run complete CI pipeline one final time
   - Confirm zero errors, warnings, or test failures
   - Document any remaining issues that require architectural changes

## Focus Areas

- **Code Quality**: Ensure all Go code follows project standards
- **Test Coverage**: Maintain passing shell integration tests
- **Security**: Address any security vulnerabilities found
- **Performance**: Fix complexity and cognitive complexity issues
- **Dependencies**: Resolve any dependency conflicts

## Success Criteria

- [ ] `just ci` completes with exit code 0
- [ ] All linting tools pass without warnings
- [ ] All shell integration tests pass
- [ ] No security vulnerabilities detected
- [ ] Code complexity within acceptable limits

Execute this systematically, fixing issues in batches and validating after each batch to avoid introducing new problems.
