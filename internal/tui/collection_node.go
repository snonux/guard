package tui

import (
	"github.com/florianbuetow/guard/internal/manager"
)

// CollectionNode represents a collection in the tree
type CollectionNode struct {
	Name         string
	FileCount    int
	IsEmpty      bool
	GuardState   GuardState
	EquivalentTo string // Name of an equivalent collection (same files), if any
	Depth        int    // Hierarchy depth for indentation
	Parent       string // Name of parent collection
	IsLastChild  bool
	AncestorLast []bool
}

// BuildCollectionNodes builds a list of collection nodes from the manager
func BuildCollectionNodes(mgr *manager.Manager) []CollectionNode {
	if mgr == nil {
		return nil
	}

	infos := BuildCollectionHierarchy(mgr)
	if len(infos) == 0 {
		return nil
	}

	nodes := make([]CollectionNode, len(infos))
	for i, info := range infos {
		// Build ancestor tracking
		var ancestorLast []bool
		if info.Depth > 0 {
			ancestorLast = make([]bool, info.Depth)
			// TODO: Track actual ancestor last status
		}

		// Determine if this is the last child at its level
		isLast := true
		if i < len(infos)-1 {
			nextInfo := infos[i+1]
			// If next item is at same or lower depth with same parent, we're not last
			if nextInfo.Depth == info.Depth && nextInfo.Parent == info.Parent {
				isLast = false
			}
		}

		nodes[i] = CollectionNode{
			Name:         info.Name,
			FileCount:    info.FileCount,
			IsEmpty:      info.IsEmpty,
			GuardState:   info.GuardState,
			EquivalentTo: info.EquivalentTo,
			Depth:        info.Depth,
			Parent:       info.Parent,
			IsLastChild:  isLast,
			AncestorLast: ancestorLast,
		}
	}

	// Second pass: fix ancestor tracking based on actual tree structure
	fixAncestorTracking(nodes)

	return nodes
}

// fixAncestorTracking fixes the ancestor last child tracking
func fixAncestorTracking(nodes []CollectionNode) {
	// For each node, determine if its ancestors are last children
	for i := range nodes {
		node := &nodes[i]
		if node.Depth == 0 {
			continue
		}

		node.AncestorLast = make([]bool, node.Depth)

		// Walk up the tree
		parentName := node.Parent
		depth := node.Depth - 1
		for depth >= 0 && parentName != "" {
			// Find the parent
			for j := i - 1; j >= 0; j-- {
				if nodes[j].Name == parentName {
					node.AncestorLast[depth] = nodes[j].IsLastChild
					parentName = nodes[j].Parent
					break
				}
			}
			depth--
		}
	}
}

// UpdateCollectionGuardStates updates the guard states for all collection nodes
func UpdateCollectionGuardStates(nodes []CollectionNode, mgr *manager.Manager) {
	if mgr == nil {
		return
	}

	for i := range nodes {
		nodes[i].GuardState = ComputeEffectiveCollectionGuardState(mgr, nodes[i].Name)
	}
}
