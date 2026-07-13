package internal

import (
	"strconv"
	"strings"
	"testing"
)

// A post-prepare-release changelog: the stable "## [1.2.3]" section already
// exists (with the moved Unreleased delta as its body) and the RC sections are
// listed in newest-first file order (rc.10 before rc.2 before rc.1).
const happyMultiRCInput = `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.2.3] - 2026-07-09

### Added

- Post-RC feature

## [1.2.3-rc.10] - 2026-07-08

### Fixed

- Bug in rc10

## [1.2.3-rc.2] - 2026-07-07

### Security

- CVE in rc2

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

### Changed

- Changed X

## [1.1.0] - 2026-06-01

### Added

- Old stuff

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3...HEAD
[1.2.3]: https://github.com/giantswarm/x/compare/v1.1.0...v1.2.3
[1.2.3-rc.10]: https://github.com/giantswarm/x/compare/v1.2.3-rc.2...v1.2.3-rc.10
[1.2.3-rc.2]: https://github.com/giantswarm/x/compare/v1.2.3-rc.1...v1.2.3-rc.2
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
[1.1.0]: https://github.com/giantswarm/x/releases/tag/v1.1.0
`

func Test_modifier_ensureReleaseCandidateChangelogsAggregated_happyMultiRC(t *testing.T) {
	m := Modifier{newVersion: "1.2.3"}

	got, err := m.ensureReleaseCandidateChangelogsAggregated([]byte(happyMultiRCInput))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	out := string(got)

	// The whole aggregated stable section, pinned including the boundary to the
	// next (untouched) RC header. This proves: category ordering (Added, Changed,
	// Fixed, Security), Unreleased-delta-first then RC ordering within a category
	// (Post-RC feature before Feature A), the note as the first Changed bullet,
	// and the plural "rc.1 through rc.10" wording (integer sort, not lexical).
	wantSection := `## [1.2.3] - 2026-07-09

### Added

- Post-RC feature
- Feature A

### Changed

- This release aggregates all changes from release candidates 1.2.3-rc.1 through 1.2.3-rc.10.
- Changed X

### Fixed

- Bug in rc10

### Security

- CVE in rc2

## [1.2.3-rc.10] - 2026-07-08
`
	if !strings.Contains(out, wantSection) {
		t.Fatalf("aggregated stable section not as expected.\n--- want to contain ---\n%s\n--- got ---\n%s", wantSection, out)
	}

	// RC sections, footer links and the Unreleased heading are left untouched.
	for _, want := range []string{
		"## [1.2.3-rc.1] - 2026-07-06",
		"## [1.2.3-rc.2] - 2026-07-07",
		"## [1.2.3-rc.10] - 2026-07-08",
		"## [Unreleased]",
		"[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1",
		"[1.1.0]: https://github.com/giantswarm/x/releases/tag/v1.1.0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to preserve %q, but it was missing", want)
		}
	}
}

func Test_modifier_ensureReleaseCandidateChangelogsAggregated_emptyDeltaSingleRC(t *testing.T) {
	input := `# Changelog

## [Unreleased]

## [1.2.3] - 2026-07-09

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

### Fixed

- Bug B

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3...HEAD
[1.2.3]: https://github.com/giantswarm/x/compare/v1.1.0...v1.2.3
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
`

	// Empty stable delta + a single RC. Exact-match output also covers the
	// singular note wording and the "Changed is always emitted with the note"
	// rule (no source has a Changed section).
	want := `# Changelog

## [Unreleased]

## [1.2.3] - 2026-07-09

### Added

- Feature A

### Changed

- This release aggregates all changes from release candidate 1.2.3-rc.1.

### Fixed

- Bug B

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

### Fixed

- Bug B

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3...HEAD
[1.2.3]: https://github.com/giantswarm/x/compare/v1.1.0...v1.2.3
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
`

	m := Modifier{newVersion: "1.2.3"}
	got, err := m.ensureReleaseCandidateChangelogsAggregated([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(got) != want {
		t.Fatalf("output mismatch.\n--- want ---\n%s\n--- got ---\n%s", want, string(got))
	}
}

func Test_modifier_ensureReleaseCandidateChangelogsAggregated_idempotent(t *testing.T) {
	m := Modifier{newVersion: "1.2.3"}

	first, err := m.ensureReleaseCandidateChangelogsAggregated([]byte(happyMultiRCInput))
	if err != nil {
		t.Fatalf("unexpected error on first run: %s", err)
	}

	second, err := m.ensureReleaseCandidateChangelogsAggregated(first)
	if err != nil {
		t.Fatalf("unexpected error on second run: %s", err)
	}

	if string(second) != string(first) {
		t.Fatalf("aggregation is not idempotent.\n--- first ---\n%s\n--- second ---\n%s", string(first), string(second))
	}
}

func Test_modifier_ensureReleaseCandidateChangelogsAggregated_noOp(t *testing.T) {
	testCases := []struct {
		name       string
		newVersion string
		input      string
	}{
		{
			name:       "case 0: RC target",
			newVersion: "1.2.3-rc.2",
			input:      happyMultiRCInput,
		},
		{
			name:       "case 1: stable target with no matching RC sections",
			newVersion: "2.0.0",
			input: `# Changelog

## [Unreleased]

## [2.0.0] - 2026-07-09

### Added

- Thing

[Unreleased]: https://github.com/giantswarm/x/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/giantswarm/x/compare/v1.9.0...v2.0.0
`,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			m := Modifier{newVersion: tc.newVersion}
			got, err := m.ensureReleaseCandidateChangelogsAggregated([]byte(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if string(got) != tc.input {
				t.Fatalf("expected no-op (unchanged content), got:\n%s", string(got))
			}
		})
	}
}

func Test_modifier_ensureReleaseCandidateChangelogsAggregated_errors(t *testing.T) {
	testCases := []struct {
		name       string
		newVersion string
		input      string
		wantInMsg  string
		matches    func(error) bool
	}{
		{
			name:       "case 0: non-canonical heading in the (ungated) stable delta",
			newVersion: "1.2.3",
			input: `# Changelog

## [Unreleased]

## [1.2.3] - 2026-07-09

### Added

- Real feature

### Testing

- Not a Keep a Changelog category

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3...HEAD
[1.2.3]: https://github.com/giantswarm/x/compare/v1.1.0...v1.2.3
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
`,
			wantInMsg: "### Testing",
			matches:   IsNonCanonicalHeading,
		},
		{
			name:       "case 1: missing stable header while RC sections exist",
			newVersion: "1.2.3",
			input: `# Changelog

## [Unreleased]

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3-rc.1...HEAD
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
`,
			wantInMsg: "## [1.2.3]",
			matches:   IsMissingStableSection,
		},
		{
			name:       "case 2: content before the first heading in the (ungated) stable delta",
			newVersion: "1.2.3",
			input: `# Changelog

## [Unreleased]

## [1.2.3] - 2026-07-09

- Orphan bullet with no category heading above it

### Added

- Real feature

## [1.2.3-rc.1] - 2026-07-06

### Added

- Feature A

[Unreleased]: https://github.com/giantswarm/x/compare/v1.2.3...HEAD
[1.2.3]: https://github.com/giantswarm/x/compare/v1.1.0...v1.2.3
[1.2.3-rc.1]: https://github.com/giantswarm/x/releases/tag/v1.2.3-rc.1
`,
			wantInMsg: "Orphan bullet with no category heading above it",
			matches:   IsOrphanContent,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			m := Modifier{newVersion: tc.newVersion}
			_, err := m.ensureReleaseCandidateChangelogsAggregated([]byte(tc.input))
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !tc.matches(err) {
				t.Fatalf("error type mismatch: %s", err)
			}
			if !strings.Contains(err.Error(), tc.wantInMsg) {
				t.Fatalf("expected error message to contain %q, got: %s", tc.wantInMsg, err.Error())
			}
		})
	}
}

func Test_changelogDocument_rcSections(t *testing.T) {
	input := `# Changelog

## [1.2.3] - 2026-07-09

## [1.2.3-rc.10] - d

## [1.2.3-rc.2] - d

## [1.2.3-rc.1] - d

## [1.2.30-rc.1] - d

## [1.1.0] - d
`
	doc := parseChangelogDocument(input)
	got := doc.rcSections("1.2.3")

	var keys []string
	for _, s := range got {
		keys = append(keys, s.versionKey)
	}
	want := []string{"1.2.3-rc.1", "1.2.3-rc.2", "1.2.3-rc.10"}
	if strings.Join(keys, ",") != strings.Join(want, ",") {
		t.Fatalf("rcSections ordering/filtering wrong: want %v, got %v (a different core like 1.2.30-rc.1 must be excluded)", want, keys)
	}
}

func Test_mergeCategorized(t *testing.T) {
	// Sources deliberately out of canonical order; output must be canonical.
	sources := [][]string{
		{"### Fixed", "- f1"},
		{"### Added", "- a1"},
	}
	got := mergeCategorized(sources, "NOTE")
	want := "### Added\n\n- a1\n\n### Changed\n\n- NOTE\n\n### Fixed\n\n- f1\n\n"
	if got != want {
		t.Fatalf("mergeCategorized mismatch.\n--- want ---\n%q\n--- got ---\n%q", want, got)
	}
}
