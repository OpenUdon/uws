package versions

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"golang.org/x/mod/module"
)

const uwsModulePath = "github.com/OpenUdon/uws"

//go:embed *.json
var embeddedSchemas embed.FS

// PathForVersion returns the best local schema path for a UWS document version.
func PathForVersion(anchorDir, version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		version = "1.0.0"
	}
	return pathForSchemaName(anchorDir, version+".json")
}

// PathForRuntimeSupplement returns the best local schema path for a runtime supplement profile.
func PathForRuntimeSupplement(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, runtimeSupplementSchemaName(profile))
}

// PathForAnsibleModuleCallSupplement returns the best local schema path for an
// Ansible module-call supplement profile.
func PathForAnsibleModuleCallSupplement(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, ansibleModuleCallSupplementSchemaName(profile))
}

// PathForAnsibleArgspec returns the best local schema path for an Ansible
// argspec document.
func PathForAnsibleArgspec(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "ansible", "1.0"))
}

// PathForAnsibleSourceProfile returns the best local schema path for an
// Ansible argspec document.
//
// Deprecated: use PathForAnsibleArgspec. Ansible argspecs are not source
// profiles as of UWS 1.7.
func PathForAnsibleSourceProfile(anchorDir, profile string) string {
	return PathForAnsibleArgspec(anchorDir, profile)
}

// PathForBrowserSourceProfile returns the best local schema path for a browser
// source profile.
func PathForBrowserSourceProfile(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "browser", "1.5"))
}

func runtimeSupplementSchemaName(profile string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		return "runtime.1.0.json"
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, "runtime.") {
		profile = "runtime." + profile
	}
	return profile + ".json"
}

func ansibleModuleCallSupplementSchemaName(profile string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		return "ansible-module-call.1.0.json"
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, "ansible-module-call.") {
		profile = "ansible-module-call." + profile
	}
	return profile + ".json"
}

func familySchemaName(profile, name, defaultVersion string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		profile = defaultVersion
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, name+".") {
		profile = name + "." + profile
	}
	return profile + ".json"
}

func pathForSchemaName(anchorDir, name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." {
		name = "1.0.0.json"
	}
	if dir := strings.TrimSpace(os.Getenv("UWS_SCHEMA_DIR")); dir != "" {
		return filepath.Join(dir, name)
	}
	if dir := strings.TrimSpace(os.Getenv("OPENUDON_UWS_SCHEMA_DIR")); dir != "" {
		return filepath.Join(dir, name)
	}
	if path, ok := packageSchemaPath(name); ok {
		return path
	}
	if path, ok := moduleCacheSchemaPath(name); ok {
		return path
	}
	if path, ok := embeddedSchemaPath(name); ok {
		return path
	}
	return filepath.Join(anchorDir, "..", "uws", "versions", name)
}

func packageSchemaPath(name string) (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", false
	}
	path := filepath.Join(filepath.Dir(file), name)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

func moduleCacheSchemaPath(name string) (string, bool) {
	version, ok := uwsModuleVersion()
	if !ok {
		return "", false
	}
	path, err := escapedModuleCachePath(uwsModulePath, version)
	if err != nil {
		return "", false
	}
	schema := filepath.Join(path, "versions", name)
	if _, err := os.Stat(schema); err == nil {
		return schema, true
	}
	return "", false
}

func uwsModuleVersion() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, dep := range info.Deps {
		if dep.Path != uwsModulePath {
			continue
		}
		if dep.Version != "" {
			return dep.Version, true
		}
		if dep.Replace != nil && dep.Replace.Version != "" {
			return dep.Replace.Version, true
		}
	}
	return "", false
}

func escapedModuleCachePath(path, version string) (string, error) {
	escapedPath, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	escapedVersion, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}
	return filepath.Join(moduleCacheDir(), escapedPath+"@"+escapedVersion), nil
}

func moduleCacheDir() string {
	if dir := strings.TrimSpace(os.Getenv("GOMODCACHE")); dir != "" {
		return dir
	}
	gopath := strings.TrimSpace(os.Getenv("GOPATH"))
	if gopath == "" {
		if home, err := os.UserHomeDir(); err == nil {
			gopath = filepath.Join(home, "go")
		}
	}
	if gopath == "" {
		return ""
	}
	first := filepath.SplitList(gopath)[0]
	if first == "" {
		return ""
	}
	return filepath.Join(first, "pkg", "mod")
}

func embeddedSchemaPath(name string) (string, bool) {
	data, err := embeddedSchemas.ReadFile(filepath.ToSlash(name))
	if err != nil {
		return "", false
	}
	sum := sha256.Sum256(data)
	dir := filepath.Join(os.TempDir(), "uws-schema", fmt.Sprintf("%x", sum[:8]))
	path := filepath.Join(dir, name)
	if existing, err := os.ReadFile(path); err == nil && bytes.Equal(existing, data) {
		return path, true
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", false
	}
	tmp, err := os.CreateTemp(dir, name+".*.tmp")
	if err != nil {
		return "", false
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return "", false
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	return path, true
}
