# Implementation Quality Check

Before you claim anything works: 

1) **Run the actual commands and verify output matches tests exactly** - Don't assume, test it
2) **Search for 'not implemented' and 'TODO' - fix ALL of them** - No stubs left behind  
3) **Ensure you moved working code, didn't rewrite it** - Preserve behavior during refactoring
4) **Pick ONE approach for each system (warnings, etc.) and remove the others** - No parallel incomplete systems
5) **Test that the core functionality actually works end-to-end** - Show me the test results, not just "it should work"

## Specific Checks:

- `grep -r "not implemented" .` - Should return nothing
- `grep -r "TODO" .` - Should return nothing or be intentional
- `go build && ./build/guard [test commands]` - Should work and match expected output
- Check for duplicate/unused systems (warning collectors, helper functions, etc.)
- Verify refactored code preserves original behavior

❌ DO NOT claim success until ALL checks pass
❌ DO NOT create beautiful structure with hollow implementation  
❌ DO NOT optimize for appearance over actual functionality

Show me actual test results and grep outputs, not assumptions.
