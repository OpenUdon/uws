// Command generateversionarchive creates the deterministic schema archive
// embedded by package schemas. The versions directory remains the source of
// truth and contains documents only.
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	source := flag.String("source", "../versions", "directory containing version documents")
	output := flag.String("output", "version-documents.zip", "archive output path")
	flag.Parse()
	if err := run(*source, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(source, output string) error {
	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("read source directory: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	var encoded bytes.Buffer
	archive := zip.NewWriter(&encoded)
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(source, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		header := &zip.FileHeader{Name: name, Method: zip.Deflate}
		header.SetMode(fs.FileMode(0o644))
		header.SetModTime(time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC))
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("create %s: %w", name, err)
		}
		if _, err := writer.Write(data); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	if err := archive.Close(); err != nil {
		return fmt.Errorf("close archive: %w", err)
	}
	if len(names) == 0 {
		return fmt.Errorf("source directory contains no JSON documents")
	}
	if err := os.WriteFile(output, encoded.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write archive: %w", err)
	}
	return nil
}
