package csm

import (
	"strings"
	"testing"
)

// Valve's stock CS2 gameinfo.gi uses tabs for indentation. This is a minimized
// but faithful shape of the real file as of late 2024 / 2025 (multiple Game
// lines inside a single SearchPaths block nested inside FileSystem).
const stockTabGameInfo = `"GameInfo"
{
	game		"Counter-Strike 2"
	title		"Counter-Strike 2"
	LayeredOnMod	csgo_imported
	FileSystem
	{
		SteamAppId		730
		SearchPaths
		{
			Game_LowViolence	csgo_lv
			Game			csgo
			Game			csgo_imported
			Game			csgo_core
			Game			core
		}
	}
}
`

// Spaces-indented variant (rare, but admins sometimes reflow the file).
const spaceIndentGameInfo = `"GameInfo"
{
    game        "Counter-Strike 2"
    FileSystem
    {
        SearchPaths
        {
            Game_LowViolence csgo_lv
            Game             csgo
            Game             csgo_imported
        }
    }
}
`

// A realistic red herring: a stray `SearchPaths` keyword appearing inside a
// comment or an unrelated sub-block. The fix MUST target the FileSystem one.
const redHerringGameInfo = `"GameInfo"
{
	game		"Counter-Strike 2"
	// Historical: some addons used to declare their own SearchPaths block here.
	SomeOtherSection
	{
		SearchPaths
		{
			Game	unrelated
		}
	}
	FileSystem
	{
		SearchPaths
		{
			Game_LowViolence	csgo_lv
			Game			csgo
		}
	}
}
`

func TestEnableMetamodInGameInfo_StockTabs(t *testing.T) {
	got, changed, warn := enableMetamodInGameInfo(stockTabGameInfo)
	if warn != "" {
		t.Fatalf("unexpected warning: %s", warn)
	}
	if !changed {
		t.Fatalf("expected changed=true on clean gameinfo.gi")
	}
	if !strings.Contains(got, "Game\tcsgo/addons/metamod") {
		t.Fatalf("expected a tab-separated Metamod entry, got:\n%s", got)
	}
	assertMetamodAfterLowViolence(t, got)
	assertMetamodIndentMatchesGameCsgo(t, got)
}

func TestEnableMetamodInGameInfo_Spaces(t *testing.T) {
	got, changed, warn := enableMetamodInGameInfo(spaceIndentGameInfo)
	if warn != "" {
		t.Fatalf("unexpected warning: %s", warn)
	}
	if !changed {
		t.Fatalf("expected changed=true on clean space-indented gameinfo.gi")
	}
	assertMetamodAfterLowViolence(t, got)
	// Inserted indent should match the other Game lines' leading spaces
	// (12 spaces in the fixture), not a tab.
	var metamodLine string
	for _, l := range strings.Split(got, "\n") {
		if strings.Contains(l, "csgo/addons/metamod") {
			metamodLine = l
			break
		}
	}
	if !strings.HasPrefix(metamodLine, "            ") {
		t.Fatalf("expected 12-space indent on Metamod line, got: %q", metamodLine)
	}
}

func TestEnableMetamodInGameInfo_Idempotent(t *testing.T) {
	once, _, _ := enableMetamodInGameInfo(stockTabGameInfo)
	twice, changed, warn := enableMetamodInGameInfo(once)
	if warn != "" {
		t.Fatalf("unexpected warning on second run: %s", warn)
	}
	if changed {
		t.Fatalf("second run should be a no-op; content diff:\n%s", twice)
	}
	if once != twice {
		t.Fatalf("idempotent run should return identical content")
	}
	if strings.Count(twice, "csgo/addons/metamod") != 1 {
		t.Fatalf("expected exactly one metamod line after repeated runs, got:\n%s", twice)
	}
}

func TestEnableMetamodInGameInfo_RedHerringSearchPaths(t *testing.T) {
	got, changed, _ := enableMetamodInGameInfo(redHerringGameInfo)
	if !changed {
		t.Fatalf("expected changed=true when a non-FileSystem SearchPaths exists")
	}
	// The unrelated SearchPaths block above must NOT have gotten the entry.
	before, _, _ := strings.Cut(got, "FileSystem")
	if strings.Contains(before, "csgo/addons/metamod") {
		t.Fatalf("metamod entry leaked into the unrelated SearchPaths block:\n%s", got)
	}
	// The FileSystem SearchPaths MUST have the entry, after Game_LowViolence.
	assertMetamodAfterLowViolence(t, got)
}

func TestEnableMetamodInGameInfo_MalformedReturnsWarning(t *testing.T) {
	_, changed, warn := enableMetamodInGameInfo(`"GameInfo" { game "CS2" }`)
	if changed {
		t.Fatalf("expected no change on malformed gameinfo")
	}
	if warn == "" {
		t.Fatalf("expected a warning for missing FileSystem block")
	}
}

func TestDisableMetamodInGameInfo_RemovesLineOnly(t *testing.T) {
	enabled, _, _ := enableMetamodInGameInfo(stockTabGameInfo)
	disabled := disableMetamodInGameInfo(enabled)
	if strings.Contains(disabled, "csgo/addons/metamod") {
		t.Fatalf("disable should strip all metamod references, got:\n%s", disabled)
	}
	// The rest of the file must be intact (byte-for-byte match with original).
	if disabled != stockTabGameInfo {
		t.Fatalf("disable should round-trip back to the original file\nwant:\n%s\ngot:\n%s", stockTabGameInfo, disabled)
	}
}

// assertMetamodAfterLowViolence checks that the metamod line appears on the
// line immediately following Game_LowViolence within the SearchPaths block.
func assertMetamodAfterLowViolence(t *testing.T, content string) {
	t.Helper()
	lines := strings.Split(content, "\n")
	for i, l := range lines {
		if strings.Contains(l, "Game_LowViolence") {
			if i+1 >= len(lines) || !strings.Contains(lines[i+1], "csgo/addons/metamod") {
				t.Fatalf("expected metamod entry on line after Game_LowViolence; lines around index %d:\n%s", i, strings.Join(lines[max0(i-1):min(i+3, len(lines))], "\n"))
			}
			return
		}
	}
	t.Fatalf("fixture missing Game_LowViolence line")
}

// assertMetamodIndentMatchesGameCsgo ensures that the inserted line shares the
// leading whitespace of the existing `Game csgo` line so it visually aligns.
func assertMetamodIndentMatchesGameCsgo(t *testing.T, content string) {
	t.Helper()
	var mmIndent, gameCsgoIndent string
	for _, l := range strings.Split(content, "\n") {
		if mmIndent == "" && strings.Contains(l, "csgo/addons/metamod") {
			mmIndent = leadingIndent(l)
		}
		if gameCsgoIndent == "" && strings.Contains(l, "Game") &&
			strings.Contains(l, "csgo") &&
			!strings.Contains(l, "addons") &&
			!strings.Contains(l, "csgo_") &&
			!strings.Contains(l, "Game_LowViolence") {
			gameCsgoIndent = leadingIndent(l)
		}
	}
	if mmIndent == "" || gameCsgoIndent == "" {
		t.Fatalf("could not locate both indents (mm=%q, gameCsgo=%q)", mmIndent, gameCsgoIndent)
	}
	if mmIndent != gameCsgoIndent {
		t.Fatalf("metamod indent %q does not match Game csgo indent %q", mmIndent, gameCsgoIndent)
	}
}

func max0(a int) int {
	if a < 0 {
		return 0
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
