package csm

import "testing"

// TestIsMetamodLinuxAsset locks in the asset-name classifier used by
// resolveMetamodURL so a cosmetic rename on AlliedModders' side can't
// silently cause CSM to pick the wrong tarball (e.g. Windows or SDK) off a
// GitHub release.
func TestIsMetamodLinuxAsset(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		want bool
	}{
		// Current AlliedModders naming for the 2.0 branch.
		{"mmsource-2.0.0-git1396-linux.tar.gz", true},
		// Older naming variants observed on releases.
		{"mmsource-1.12.0-git1234-linux.tar.gz", true},
		// Architecture suffixes should still match.
		{"mmsource-2.0.0-git1400-linux-x64.tar.gz", true},
		{"MMSource-2.0.0-git1401-Linux.tar.gz", true},

		// Negative cases: wrong OS.
		{"mmsource-2.0.0-git1396-windows.zip", false},
		{"mmsource-2.0.0-git1396-mac.tar.gz", false},
		// Negative cases: wrong archive format.
		{"mmsource-2.0.0-git1396-linux.zip", false},
		{"mmsource-2.0.0-git1396-linux.tar.xz", false},
		// Negative cases: SDK/source-only artifacts occasionally attached.
		{"mmsource-2.0.0-sdk-linux.tar.gz", false},
		{"source-sdk-linux.tar.gz", false},
		// Negative cases: empty or obviously unrelated.
		{"", false},
		{"README.md", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := isMetamodLinuxAsset(tc.name); got != tc.want {
				t.Fatalf("isMetamodLinuxAsset(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
