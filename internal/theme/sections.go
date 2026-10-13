package theme

import (
	"strconv"
	"strings"
)

// Section keys for the fixed parts of the public page. Custom content blocks
// join the same ordering as "block:<id>".
const (
	SectionHeader = "header"
	SectionSearch = "search"
	SectionNav    = "nav"
	SectionMenu   = "menu"
	SectionHours  = "hours"
	SectionFooter = "footer"
)

// FixedSections are the built-in sections, in their default order.
var FixedSections = []string{
	SectionHeader, SectionSearch, SectionNav, SectionMenu, SectionHours, SectionFooter,
}

// blockPrefix marks a section entry that refers to a custom block.
const blockPrefix = "block:"

// IsBlockSection reports whether a section key points at a custom block, and
// returns that block's id.
func IsBlockSection(key string) (int64, bool) {
	if !strings.HasPrefix(key, blockPrefix) {
		return 0, false
	}
	id, err := strconv.ParseInt(strings.TrimPrefix(key, blockPrefix), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// BlockSectionKey builds the section key for a block id.
func BlockSectionKey(id int64) string { return blockPrefix + strconv.FormatInt(id, 10) }

// DefaultSections is the starting page order for a new venue.
func DefaultSections() []string {
	out := make([]string, len(FixedSections))
	copy(out, FixedSections)
	return out
}

// ReconcileSections merges the owner's saved order with reality.
//
// The stored order can drift: a block gets deleted, a new one is added, or a
// future release introduces another fixed section. This drops entries that no
// longer exist, keeps the order of those that do, and appends anything missing
// just before the footer so new content is visible rather than silently lost.
func ReconcileSections(stored []string, blockIDs []int64) []string {
	known := make(map[string]bool, len(FixedSections)+len(blockIDs))
	for _, s := range FixedSections {
		known[s] = true
	}
	for _, id := range blockIDs {
		known[BlockSectionKey(id)] = true
	}

	out := make([]string, 0, len(known))
	seen := make(map[string]bool, len(known))

	for _, s := range stored {
		if known[s] && !seen[s] {
			out = append(out, s)
			seen[s] = true
		}
	}

	// Anything not mentioned in the stored order goes in before the footer.
	var missing []string
	for _, s := range FixedSections {
		if !seen[s] {
			missing = append(missing, s)
		}
	}
	for _, id := range blockIDs {
		key := BlockSectionKey(id)
		if !seen[key] {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		return out
	}

	// New content goes at the end, except that it should not land after a
	// trailing footer. Anchoring on a *trailing* footer rather than wherever
	// the footer happens to sit matters: an owner who deliberately moved the
	// footer to the top would otherwise see new sections pushed above it.
	insertAt := len(out)
	if len(out) > 0 && out[len(out)-1] == SectionFooter {
		insertAt = len(out) - 1
	}

	merged := make([]string, 0, len(out)+len(missing))
	merged = append(merged, out[:insertAt]...)
	merged = append(merged, missing...)
	merged = append(merged, out[insertAt:]...)
	return merged
}

// SanitizeSections keeps only well-formed entries, dropping duplicates. It runs
// during Config.Sanitize, before the block list is known; ReconcileSections
// does the content-aware pass later.
func SanitizeSections(in []string) []string {
	if len(in) == 0 {
		return DefaultSections()
	}

	fixed := make(map[string]bool, len(FixedSections))
	for _, s := range FixedSections {
		fixed[s] = true
	}

	out := make([]string, 0, len(in))
	seen := make(map[string]bool, len(in))

	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		_, isBlock := IsBlockSection(s)
		if !fixed[s] && !isBlock {
			continue
		}
		out = append(out, s)
		seen[s] = true
	}

	if len(out) == 0 {
		return DefaultSections()
	}
	return out
}
