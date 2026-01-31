package tui

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
)

// StringWidth returns the display width of a string, accounting for wide characters
// and ANSI escape sequences (which are not counted as visible characters)
func StringWidth(s string) int {
	return ansi.PrintableRuneWidth(s)
}

// TruncateMiddle truncates a string to fit within maxWidth, adding "..." in the middle
// if truncation is needed. This preserves both the beginning and end of the string.
func TruncateMiddle(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	width := StringWidth(s)
	if width <= maxWidth {
		return s
	}

	// Need at least 4 chars for "a...b"
	if maxWidth < 5 {
		// Just truncate from the right
		return TruncateRight(s, maxWidth)
	}

	ellipsis := "..."
	ellipsisWidth := 3

	// Calculate how many chars we can show on each side
	availableWidth := maxWidth - ellipsisWidth
	leftWidth := availableWidth / 2
	rightWidth := availableWidth - leftWidth

	// Get left part
	left := truncateToWidth(s, leftWidth)

	// Get right part (from the end)
	right := truncateFromEnd(s, rightWidth)

	return left + ellipsis + right
}

// TruncateRight truncates a string from the right to fit within maxWidth
func TruncateRight(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	width := StringWidth(s)
	if width <= maxWidth {
		return s
	}

	return truncateToWidth(s, maxWidth)
}

// truncateToWidth returns a prefix of s that fits within maxWidth
func truncateToWidth(s string, maxWidth int) string {
	var result []rune
	currentWidth := 0

	for _, r := range s {
		charWidth := runewidth.RuneWidth(r)
		if currentWidth+charWidth > maxWidth {
			break
		}
		result = append(result, r)
		currentWidth += charWidth
	}

	return string(result)
}

// truncateFromEnd returns a suffix of s that fits within maxWidth
func truncateFromEnd(s string, maxWidth int) string {
	runes := []rune(s)
	var result []rune
	currentWidth := 0

	// Walk backwards from the end
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		charWidth := runewidth.RuneWidth(r)
		if currentWidth+charWidth > maxWidth {
			break
		}
		result = append([]rune{r}, result...)
		currentWidth += charWidth
	}

	return string(result)
}

// PadRight pads a string with spaces to reach the desired width
func PadRight(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

// PadLeft pads a string with spaces on the left to reach the desired width
func PadLeft(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := width - currentWidth
	return strings.Repeat(" ", padding) + s
}
