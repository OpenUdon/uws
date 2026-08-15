package schemas

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionsDirectoryContainsDocumentsOnly(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "versions"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".go") {
			t.Fatalf("versions directory contains Go source %q", entry.Name())
		}
	}
}

func TestEmbeddedVersionDocumentsMatchSource(t *testing.T) {
	archive, err := zip.NewReader(bytes.NewReader(embeddedVersionDocuments), int64(len(embeddedVersionDocuments)))
	if err != nil {
		t.Fatal(err)
	}
	embedded := make(map[string]bool, len(archive.File))
	for _, file := range archive.File {
		embedded[file.Name] = true
	}
	entries, err := os.ReadDir(filepath.Join("..", "versions"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		name := entry.Name()
		want, err := os.ReadFile(filepath.Join("..", "versions", name))
		if err != nil {
			t.Fatal(err)
		}
		got, err := embeddedSchemaDocument(name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("embedded %s is stale; run go generate ./schemas", name)
		}
		delete(embedded, name)
	}
	if len(embedded) != 0 {
		t.Fatalf("embedded archive contains documents absent from versions/: %v", embedded)
	}
}

func TestPathForVersionFindsReadableSchema(t *testing.T) {
	for _, version := range []string{"", "1.0.0", "1.1.0", "1.1.1", "1.2.0", "1.3.0", "1.4.0", "1.5.0", "1.6.0", "1.7.0"} {
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

func TestPathForBrowserSourceProfileFindsReadableSchema(t *testing.T) {
	for _, profile := range []string{"", "1.5", "1.5.json", "browser.1.5", "browser.1.5.json", "uws.browser.1.5", "uws.browser.1.5.json"} {
		path := PathForBrowserSourceProfile(t.TempDir(), profile)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("browser source profile path for profile %q is not readable: %s: %v", profile, path, err)
		}
	}
}

func TestProfilePathsHonorSchemaDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UWS_SCHEMA_DIR", dir)
	path := PathForBrowserSourceProfile(".", "uws.browser.1.5.json")
	want := filepath.Join(dir, "browser.1.5.json")
	if path != want {
		t.Fatalf("browser schema path = %q, want %q", path, want)
	}
}

func TestEmbeddedSchemaPathFindsReadableSchema(t *testing.T) {
	for _, name := range []string{"1.0.0.json", "1.1.0.json", "1.1.1.json", "1.2.0.json", "1.3.0.json", "1.4.0.json", "1.5.0.json", "1.6.0.json", "1.7.0.json", "runtime.1.0.json", "browser.1.5.json", "browser-authentication.1.0.json", "browser-authentication-call.1.0.json"} {
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
