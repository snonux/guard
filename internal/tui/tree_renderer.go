package tui

import "strings"

// Tree drawing characters
const (
	TreeBranch     = "├─ "
	TreeLastBranch = "└─ "
	TreeVertical   = "│  "
	TreeEmpty      = "   "
)

// Folder expansion indicators
const (
	FolderExpanded  = "▼ "
	FolderCollapsed = "▶ "
	FileIndicator   = "  " // No indicator for files, just spacing
)

// BuildTreePrefix generates the prefix string for a tree node based on its depth
// and whether its ancestors are the last child in their respective levels
func BuildTreePrefix(depth int, isLastAtLevel []bool) string {
	if depth == 0 {
		return ""
	}

	var sb strings.Builder

	// Add prefix for each ancestor level
	for i := range depth - 1 {
		if i < len(isLastAtLevel) && isLastAtLevel[i] {
			sb.WriteString(TreeEmpty)
		} else {
			sb.WriteString(TreeVertical)
		}
	}

	// Add the branch for this level
	if depth > 0 && len(isLastAtLevel) > 0 && isLastAtLevel[len(isLastAtLevel)-1] {
		sb.WriteString(TreeLastBranch)
	} else {
		sb.WriteString(TreeBranch)
	}

	return sb.String()
}

// GetFolderIndicator returns the appropriate expansion indicator for a folder
func GetFolderIndicator(expanded bool) string {
	if expanded {
		return FolderExpanded
	}
	return FolderCollapsed
}

// GetFileIndicator returns the indicator for a file (spacing to align with folders)
func GetFileIndicator() string {
	return FileIndicator
}

// TreePrefixWidth returns the width of a tree prefix at a given depth
func TreePrefixWidth(depth int) int {
	if depth == 0 {
		return 0
	}
	// Each level adds 3 characters (e.g., "│  " or "├─ ")
	return depth * 3
}
