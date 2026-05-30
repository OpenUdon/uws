package versions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathForVersionFindsReadableSchema(t *testing.T) {
	for _, version := range []string{"", "1.0.0", "1.1.0", "1.1.1", "1.2.0", "1.3.0", "1.4.0", "1.5.0"} {
		path := PathForVersion(t.TempDir(), version)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("schema path for version %q is not readable: %s: %v", version, path, err)
		}
	}
}

func TestPathForVersionHonorsSchemaDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UWS_SCHEMA_DIR", dir)
	if got, want := PathForVersion(".", "1.4.0"), filepath.Join(dir, "1.4.0.json"); got != want {
		t.Fatalf("PathForVersion = %q, want %q", got, want)
	}
}

func TestPathForVersionHonorsOpenUdonAlias(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UWS_SCHEMA_DIR", "")
	t.Setenv("OPENUDON_UWS_SCHEMA_DIR", dir)
	if got, want := PathForVersion(".", "1.3.0"), filepath.Join(dir, "1.3.0.json"); got != want {
		t.Fatalf("PathForVersion = %q, want %q", got, want)
	}
}

func TestPathForRuntimeSupplementFindsReadableSchema(t *testing.T) {
	for _, profile := range []string{"", "1.0", "runtime.1.0", "uws.runtime.1.0"} {
		path := PathForRuntimeSupplement(t.TempDir(), profile)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("runtime schema path for profile %q is not readable: %s: %v", profile, path, err)
		}
	}
}

func TestEmbeddedSchemaPathFindsReadableSchema(t *testing.T) {
	for _, name := range []string{"1.0.0.json", "1.1.0.json", "1.1.1.json", "1.2.0.json", "1.3.0.json", "1.4.0.json", "1.5.0.json", "runtime.1.0.json", "browser.1.5.json"} {
		path, ok := embeddedSchemaPath(name)
		if !ok {
			t.Fatalf("embedded schema path %s not found", name)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("embedded schema path %s is not readable: %v", path, err)
		}
	}
}

func TestModuleCacheSchemaPathFindsDependencySchema(t *testing.T) {
	if _, ok := uwsModuleVersion(); !ok {
		t.Skip("uws module version is unavailable, likely because the package is under direct test or workspace-replaced")
	}
	path, ok := moduleCacheSchemaPath("1.0.0.json")
	if !ok {
		t.Fatalf("module cache schema path not found")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("module cache schema path %s is not readable: %v", path, err)
	}
}
