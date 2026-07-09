package internal

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/giantswarm/gitsemver/v2/pkg/gitsemver"
	"github.com/giantswarm/microerror"
)

// changelogCategory is a Keep a Changelog H3 section name.
type changelogCategory string

const (
	categoryAdded      changelogCategory = "Added"
	categoryChanged    changelogCategory = "Changed"
	categoryDeprecated changelogCategory = "Deprecated"
	categoryRemoved    changelogCategory = "Removed"
	categoryFixed      changelogCategory = "Fixed"
	categorySecurity   changelogCategory = "Security"
)

// canonicalCategories is the fixed emit order for the aggregated section, per
// https://keepachangelog.com/en/1.1.0/.
var canonicalCategories = []changelogCategory{
	categoryAdded, categoryChanged, categoryDeprecated,
	categoryRemoved, categoryFixed, categorySecurity,
}

// aggregationNotePrefix is the leading text of the note bullet prepended under
// "### Changed" of a promoted stable section. It doubles as the idempotency
// marker: if the stable section already contains it, aggregation is skipped.
const aggregationNotePrefix = "This release aggregates all changes from release candidate"

var (
	sectionHeaderRegex = regexp.MustCompile(`^## \[([^\]]+)\]`)
	linkRefLineRegex   = regexp.MustCompile(`^\[[^\]]+\]:\s*https?://`)
	rcVersionRegex     = regexp.MustCompile(`^(.+)-rc\.([0-9]+)$`)
)

// AggregateReleaseCandidateChangelogs merges the changelog sections of a stable
// release's release candidates into the stable section, when the new version is
// a stable promotion of an existing RC series. It is a no-op otherwise.
//
// It must run after AddReleaseToChangelogMd, which creates the "## [<version>]"
// section this method aggregates into.
func (m *Modifier) AggregateReleaseCandidateChangelogs() error {
	err := modifyFile(filepath.Join(m.workingDir, FileChangelogMd), m.aggregateReleaseCandidateChangelogs)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *Modifier) aggregateReleaseCandidateChangelogs(content []byte) ([]byte, error) {
	// Only a stable target can be a promotion of an RC series. RC and dev builds
	// carry their own changelog and are never aggregated.
	if !gitsemver.IsValidStable(m.newVersion) {
		return content, nil
	}

	doc := parseChangelogDocument(string(content))

	// Release-candidate sections for this exact core version, oldest -> newest.
	rcs := doc.rcSections(m.newVersion)
	if len(rcs) == 0 {
		// Stable release with no matching RC entries: not a promotion.
		return content, nil
	}

	stable, ok := doc.section(m.newVersion)
	if !ok {
		return nil, microerror.Maskf(missingStableSectionError,
			"changelog section %#q not found while %d release-candidate section(s) exist for it; expected AddReleaseToChangelogMd to have created it",
			"## ["+m.newVersion+"]", len(rcs))
	}

	stableBody := doc.body(stable)
	if containsAggregationNote(stableBody) {
		// Already aggregated (idempotent re-run).
		return content, nil
	}

	// Source stream: the stable section body (the "Unreleased" delta that
	// AddReleaseToChangelogMd moved into it) first, then each RC oldest -> newest.
	sources := make([][]string, 0, len(rcs)+1)
	sources = append(sources, stableBody)
	for _, rc := range rcs {
		sources = append(sources, doc.body(rc))
	}

	// RC sections already passed the changelog validator (six canonical
	// categories only), but the stable body originates from an "Unreleased"
	// delta that is not gated by that validation. The category merge below only
	// emits canonical categories, so refuse to silently drop content: fail if a
	// non-canonical H3 appears in any source.
	if heading, bad := firstNonCanonicalHeading(sources); bad {
		return nil, microerror.Maskf(nonCanonicalHeadingError,
			"changelog for %#q contains non-Keep-a-Changelog heading %#q; allowed: Added, Changed, Deprecated, Removed, Fixed, Security",
			m.newVersion, "### "+heading)
	}

	merged := mergeCategorized(sources, aggregationNote(rcs))

	return doc.spliceBody(stable, merged), nil
}

// changelogSection is one "## [...]" block. Body spans [bodyStart, bodyEnd) in
// the document's line slice.
type changelogSection struct {
	versionKey string
	headerLine int
	bodyStart  int
	bodyEnd    int
}

// changelogDocument is a line-indexed CHANGELOG.md. footerStart is the first
// line of the trailing block of blank + link-reference lines (or len(lines) if
// there is none); sections are only parsed above it, so a section body can
// never spill into the footer link definitions.
type changelogDocument struct {
	lines       []string
	sections    []changelogSection
	footerStart int
}

func parseChangelogDocument(content string) changelogDocument {
	lines := strings.Split(content, "\n")

	// Isolate the trailing footer: the maximal run of blank / link-reference
	// lines at the end of the file.
	footerStart := len(lines)
	for footerStart > 0 {
		l := lines[footerStart-1]
		if strings.TrimSpace(l) == "" || linkRefLineRegex.MatchString(l) {
			footerStart--
			continue
		}
		break
	}

	var sections []changelogSection
	for i := 0; i < footerStart; i++ {
		match := sectionHeaderRegex.FindStringSubmatch(lines[i])
		if match == nil {
			continue
		}
		if n := len(sections); n > 0 {
			sections[n-1].bodyEnd = i
		}
		sections = append(sections, changelogSection{
			versionKey: match[1],
			headerLine: i,
			bodyStart:  i + 1,
			bodyEnd:    footerStart,
		})
	}

	return changelogDocument{lines: lines, sections: sections, footerStart: footerStart}
}

func (d changelogDocument) section(versionKey string) (changelogSection, bool) {
	for _, s := range d.sections {
		if s.versionKey == versionKey {
			return s, true
		}
	}
	return changelogSection{}, false
}

// rcSections returns the "## [<core>-rc.N]" sections, sorted by the integer N
// (so rc.1 < rc.2 < rc.10, not lexicographically).
func (d changelogDocument) rcSections(core string) []changelogSection {
	type numbered struct {
		section changelogSection
		n       int
	}

	var matched []numbered
	for _, s := range d.sections {
		match := rcVersionRegex.FindStringSubmatch(s.versionKey)
		if match == nil || match[1] != core {
			continue
		}
		n, err := strconv.Atoi(match[2])
		if err != nil {
			continue
		}
		matched = append(matched, numbered{section: s, n: n})
	}

	sort.Slice(matched, func(i, j int) bool { return matched[i].n < matched[j].n })

	out := make([]changelogSection, len(matched))
	for i, ns := range matched {
		out[i] = ns.section
	}
	return out
}

// body returns a section's content lines with link-reference definitions
// removed and leading/trailing blank lines trimmed.
func (d changelogDocument) body(s changelogSection) []string {
	var out []string
	for _, line := range d.lines[s.bodyStart:s.bodyEnd] {
		if linkRefLineRegex.MatchString(line) {
			continue
		}
		out = append(out, line)
	}
	for len(out) > 0 && strings.TrimSpace(out[0]) == "" {
		out = out[1:]
	}
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
		out = out[:len(out)-1]
	}
	return out
}

// spliceBody replaces the body of section s with merged (a block ending in a
// newline), leaving every other line - preamble, "## [Unreleased]", older
// sections, RC sections and the footer - untouched.
func (d changelogDocument) spliceBody(s changelogSection, merged string) []byte {
	out := make([]string, 0, len(d.lines))
	out = append(out, d.lines[:s.bodyStart]...) // up to and including the header
	out = append(out, "")                       // blank line after the header
	out = append(out, strings.Split(strings.TrimRight(merged, "\n"), "\n")...)

	tail := d.lines[s.bodyEnd:]
	// Separate the merged body from the following content with one blank line,
	// unless the tail already starts blank.
	if len(tail) > 0 && strings.TrimSpace(tail[0]) != "" {
		out = append(out, "")
	}
	out = append(out, tail...)

	return []byte(strings.Join(out, "\n"))
}

// categoryOf reports whether line is a "### <name>" heading and returns the
// trimmed name. Mirrors the awk `/^### /` + trailing-space-strip behaviour.
func categoryOf(line string) (string, bool) {
	if !strings.HasPrefix(line, "### ") {
		return "", false
	}
	return strings.TrimRightFunc(line[len("### "):], unicode.IsSpace), true
}

func isCanonicalCategory(name string) bool {
	for _, c := range canonicalCategories {
		if string(c) == name {
			return true
		}
	}
	return false
}

// firstNonCanonicalHeading returns the first "### " heading across sources whose
// name is not a canonical Keep a Changelog category.
func firstNonCanonicalHeading(sources [][]string) (string, bool) {
	for _, src := range sources {
		for _, line := range src {
			if name, ok := categoryOf(line); ok && !isCanonicalCategory(name) {
				return name, true
			}
		}
	}
	return "", false
}

// containsAggregationNote reports whether body already carries the aggregation
// note (idempotency marker).
func containsAggregationNote(body []string) bool {
	for _, line := range body {
		if strings.Contains(line, aggregationNotePrefix) {
			return true
		}
	}
	return false
}

// aggregationNote builds the note bullet text for the given RC sections
// (oldest -> newest), matching singular vs. plural phrasing.
func aggregationNote(rcs []changelogSection) string {
	first := rcs[0].versionKey
	last := rcs[len(rcs)-1].versionKey
	if first == last {
		return fmt.Sprintf("%s %s.", aggregationNotePrefix, first)
	}
	return fmt.Sprintf("%ss %s through %s.", aggregationNotePrefix, first, last)
}

// mergeCategorized buckets bullets from all sources by category (preserving
// insertion order within each category), prepends note as the first "Changed"
// bullet, and emits the categories in canonical order. Blank lines and any text
// before the first category heading are dropped.
func mergeCategorized(sources [][]string, note string) string {
	content := make(map[changelogCategory][]string)

	current := changelogCategory("")
	for _, src := range sources {
		for _, line := range src {
			if name, ok := categoryOf(line); ok {
				current = changelogCategory(name)
				continue
			}
			if current == "" || strings.TrimSpace(line) == "" {
				continue
			}
			content[current] = append(content[current], line)
		}
	}

	var b strings.Builder
	for _, c := range canonicalCategories {
		body := content[c]
		if c == categoryChanged {
			body = append([]string{"- " + note}, body...)
		}
		if len(body) == 0 {
			continue
		}
		b.WriteString("### " + string(c) + "\n\n")
		for _, line := range body {
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}
