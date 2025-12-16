package migorate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type SQLFile struct {
	Name    string
	Content string
}

func FromEmbedFS(fs embed.FS) ([]SQLFile, error) {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("reading FS failed: %w", err)
	}

	var files []SQLFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		b, err := fs.ReadFile(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("reading file failed: %s: %w", entry.Name(), err)
		}

		files = append(files, SQLFile{
			Name:    entry.Name(),
			Content: string(b),
		})
	}

	slices.SortFunc(files, func(a, b SQLFile) int {
		return strings.Compare(a.Name, b.Name)
	})

	return files, nil
}

func FromDir(dir string) ([]SQLFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading dir failed: %s: %w", dir, err)
	}

	var files []SQLFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.ToLower(filepath.Ext(entry.Name())) != ".sql" {
			continue
		}

		b, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading file failed: %s: %w", entry.Name(), err)
		}

		files = append(files, SQLFile{
			Name:    entry.Name(),
			Content: string(b),
		})
	}

	slices.SortFunc(files, func(a, b SQLFile) int {
		return strings.Compare(a.Name, b.Name)
	})

	return files, nil
}
