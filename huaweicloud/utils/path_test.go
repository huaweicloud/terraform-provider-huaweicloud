package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandHome(t *testing.T) {
	t.Setenv("HOME", "/home/testuser")
	t.Setenv("USERPROFILE", "/home/testuser")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "empty path",
			path: "",
			want: "",
		},
		{
			name: "absolute path",
			path: "/etc/hosts",
			want: "/etc/hosts",
		},
		{
			name: "relative path",
			path: "relative/path",
			want: "relative/path",
		},
		{
			name: "tilde only",
			path: "~",
			want: "/home/testuser",
		},
		{
			name: "tilde with slash",
			path: "~/foo/bar",
			want: filepath.Join("/home/testuser", "foo/bar"),
		},
		{
			name: "tilde with backslash",
			path: "~\\foo\\bar",
			want: filepath.Join("/home/testuser", "\\foo\\bar"),
		},
		{
			name:    "other user path",
			path:    "~otheruser/foo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandHome(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ExpandHome() expected error, got %q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ExpandHome() unexpected error: %s", err)
			}

			if got != tt.want {
				t.Fatalf("ExpandHome() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserHomeDir(t *testing.T) {
	t.Setenv("HOME", "/home/testuser")
	t.Setenv("USERPROFILE", "/home/testuser")

	home, err := UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() unexpected error: %s", err)
	}

	if home != "/home/testuser" {
		t.Fatalf("UserHomeDir() = %q, want %q", home, "/home/testuser")
	}
}

// legacyExpandHome replicates github.com/mitchellh/go-homedir Expand() using the
// same home resolution as legacy Dir() when only environment variables apply.
func legacyExpandHome(path string, home string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	if home == "" {
		return "", errors.New("home directory is not defined")
	}

	return filepath.Join(home, path[1:]), nil
}

// legacyDefaultSharedConfigFile replicates the historical readConfig default path
// construction before Expand was applied.
func legacyDefaultSharedConfigFile(homeEnv, userProfile string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(userProfile, ".hcloud", "config.json")
	}

	return filepath.Join(homeEnv, ".hcloud", "config.json")
}

// legacyWindowsHomeDir replicates go-homedir dirWindows() env lookup order.
func legacyWindowsHomeDir(homeEnv, userProfile, homeDrive, homePath string) (string, error) {
	if homeEnv != "" {
		return homeEnv, nil
	}

	if userProfile != "" {
		return userProfile, nil
	}

	if homeDrive != "" && homePath != "" {
		return homeDrive + homePath, nil
	}

	return "", errors.New("HOMEDRIVE, HOMEPATH, or USERPROFILE are blank")
}

func TestExpandHomeMatchesLegacyGoHomedirWhenHomeIsSet(t *testing.T) {
	home := "/home/legacy-user"
	paths := []string{
		"",
		"/etc/hosts",
		"relative/path",
		"~",
		"~/foo/bar",
		"~\\foo\\bar",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)

			got, err := ExpandHome(path)
			if err != nil {
				t.Fatalf("ExpandHome() unexpected error: %s", err)
			}

			want, err := legacyExpandHome(path, home)
			if err != nil {
				t.Fatalf("legacyExpandHome() unexpected error: %s", err)
			}

			if got != want {
				t.Fatalf("ExpandHome() = %q, legacy = %q", got, want)
			}
		})
	}
}

func TestExpandHomeMatchesLegacyGoHomedirErrorCases(t *testing.T) {
	t.Setenv("HOME", "/home/legacy-user")
	t.Setenv("USERPROFILE", "/home/legacy-user")

	_, err := ExpandHome("~otheruser/foo")
	if err == nil {
		t.Fatal("ExpandHome() expected error for other-user path")
	}

	_, err = legacyExpandHome("~otheruser/foo", "/home/legacy-user")
	if err == nil {
		t.Fatal("legacyExpandHome() expected error for other-user path")
	}
}

func TestReadConfigDefaultPathMatchesLegacyWhenHomeEnvIsSet(t *testing.T) {
	home := "/home/legacy-user"
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	legacyPath := legacyDefaultSharedConfigFile(home, home)
	legacyExpanded, err := legacyExpandHome(legacyPath, home)
	if err != nil {
		t.Fatalf("legacyExpandHome() unexpected error: %s", err)
	}

	newHome, err := UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() unexpected error: %s", err)
	}

	newPath := filepath.Join(newHome, ".hcloud", "config.json")
	newExpanded, err := ExpandHome(newPath)
	if err != nil {
		t.Fatalf("ExpandHome() unexpected error: %s", err)
	}

	if legacyExpanded != newExpanded {
		t.Fatalf("default shared config path mismatch: legacy = %q, new = %q", legacyExpanded, newExpanded)
	}
}

func TestReadConfigExplicitTildePathMatchesLegacyWhenHomeEnvIsSet(t *testing.T) {
	home := "/home/legacy-user"
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	input := "~/.hcloud/config.json"

	legacyExpanded, err := legacyExpandHome(input, home)
	if err != nil {
		t.Fatalf("legacyExpandHome() unexpected error: %s", err)
	}

	newExpanded, err := ExpandHome(input)
	if err != nil {
		t.Fatalf("ExpandHome() unexpected error: %s", err)
	}

	if legacyExpanded != newExpanded {
		t.Fatalf("explicit tilde path mismatch: legacy = %q, new = %q", legacyExpanded, newExpanded)
	}
}

// stdlibUserHomeDir mimics os.UserHomeDir() env lookup for test injection.
func stdlibUserHomeDir(goos string, getenv func(string) string) func() (string, error) {
	return func() (string, error) {
		envKey := "HOME"
		envName := "$HOME"
		if goos == "windows" {
			envKey = "USERPROFILE"
			envName = "%USERPROFILE%"
		}

		if home := getenv(envKey); home != "" {
			return home, nil
		}

		return "", errors.New(envName + " is not defined")
	}
}

func setWindowsHomeEnvs(t *testing.T, homeEnv, userProfile, homeDrive, homePath string) {
	t.Helper()
	t.Setenv("HOME", homeEnv)
	t.Setenv("USERPROFILE", userProfile)
	t.Setenv("HOMEDRIVE", homeDrive)
	t.Setenv("HOMEPATH", homePath)
}

func TestUserHomeDirMatchesLegacyWindowsEnvOrder(t *testing.T) {
	cases := []struct {
		name        string
		homeEnv     string
		userProfile string
		homeDrive   string
		homePath    string
	}{
		{
			name:        "HOME only",
			homeEnv:     `C:\custom-home`,
			userProfile: "",
			homeDrive:   "",
			homePath:    "",
		},
		{
			name:        "USERPROFILE only",
			homeEnv:     "",
			userProfile: `C:\Users\profile`,
			homeDrive:   "",
			homePath:    "",
		},
		{
			name:        "HOMEDRIVE and HOMEPATH only",
			homeEnv:     "",
			userProfile: "",
			homeDrive:   "C:",
			homePath:    `\Users\fallback`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setWindowsHomeEnvs(t, tc.homeEnv, tc.userProfile, tc.homeDrive, tc.homePath)

			got, err := resolveUserHomeDir("windows", stdlibUserHomeDir("windows", os.Getenv), os.Getenv)
			if err != nil {
				t.Fatalf("resolveUserHomeDir() unexpected error: %s", err)
			}

			want, err := legacyWindowsHomeDir(tc.homeEnv, tc.userProfile, tc.homeDrive, tc.homePath)
			if err != nil {
				t.Fatalf("legacyWindowsHomeDir() unexpected error: %s", err)
			}

			if got != want {
				t.Fatalf("resolveUserHomeDir() = %q, legacy = %q", got, want)
			}
		})
	}
}

func TestUserHomeDirDiffersFromLegacyWindowsWhenHomeAndUserProfileBothSet(t *testing.T) {
	setWindowsHomeEnvs(t, `C:\custom-home`, `C:\Users\profile`, "", "")

	got, err := resolveUserHomeDir("windows", stdlibUserHomeDir("windows", os.Getenv), os.Getenv)
	if err != nil {
		t.Fatalf("resolveUserHomeDir() unexpected error: %s", err)
	}

	legacy, err := legacyWindowsHomeDir(`C:\custom-home`, `C:\Users\profile`, "", "")
	if err != nil {
		t.Fatalf("legacyWindowsHomeDir() unexpected error: %s", err)
	}

	if got == legacy {
		t.Fatalf("expected divergence when HOME and USERPROFILE differ, both returned %q", got)
	}

	if got != `C:\Users\profile` {
		t.Fatalf("resolveUserHomeDir() = %q, want %q", got, `C:\Users\profile`)
	}

	if legacy != `C:\custom-home` {
		t.Fatalf("legacyWindowsHomeDir() = %q, want %q", legacy, `C:\custom-home`)
	}
}

func TestUserHomeDirFailsWhenHomeEnvUnsetOnUnix(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	_, err := resolveUserHomeDir("linux", stdlibUserHomeDir("linux", os.Getenv), os.Getenv)
	if err == nil {
		t.Fatal("resolveUserHomeDir() expected error when HOME is unset")
	}

	_, err = UserHomeDir()
	if runtime.GOOS == "windows" {
		t.Skip("UserHomeDir() Unix unset-HOME check runs on non-Windows hosts")
	}
	if err == nil {
		t.Fatal("UserHomeDir() expected error when HOME is unset")
	}
}
