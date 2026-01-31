package tui

// GuardState represents the guard status of a file, folder, or collection
type GuardState int

const (
	// GuardStateNotRegistered means the file is not in the registry (files only)
	GuardStateNotRegistered GuardState = iota
	// GuardStateNoCollection means the folder has no collection (folders only)
	GuardStateNoCollection
	// GuardStateUnguarded means registered but guard is off
	GuardStateUnguarded
	// GuardStateExplicit means explicitly guarded (file.Guard = true)
	GuardStateExplicit
	// GuardStateImplicit means implicitly guarded via collection (file.Guard = false but collection.Guard = true)
	GuardStateImplicit
	// GuardStateMixed means some items are guarded, some are not (folders/collections only)
	GuardStateMixed
)

// String returns the display indicator for the guard state
func (s GuardState) String() string {
	switch s {
	case GuardStateNotRegistered:
		return "[ ]"
	case GuardStateNoCollection:
		return "[ ]"
	case GuardStateUnguarded:
		return "[-]"
	case GuardStateExplicit:
		return "[G]"
	case GuardStateImplicit:
		return "[g]"
	case GuardStateMixed:
		return "[~]"
	default:
		return "[?]"
	}
}

// Panel represents which panel is currently active
type Panel int

const (
	PanelFiles Panel = iota
	PanelCollections
)

// MinTerminalWidth is the minimum terminal width required for the TUI
const MinTerminalWidth = 40

// MinTerminalHeight is the minimum terminal height required for the TUI
const MinTerminalHeight = 25
