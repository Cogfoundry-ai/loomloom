package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrOutputDirNotEmpty    = errors.New("output directory is not empty")
	ErrOutputPathNotDir     = errors.New("output path exists and is not a directory")
	ErrOutputDirUnavailable = errors.New("output directory is unavailable")
)

func checkOutputDir(outputDir string, dryRun bool) error {
	info, err := os.Stat(outputDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%w: %s", ErrOutputPathNotDir, outputDir)
		}
		entries, err := os.ReadDir(outputDir)
		if err != nil {
			return fmt.Errorf("%w: read output directory: %w", ErrOutputDirUnavailable, err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("%w: %s", ErrOutputDirNotEmpty, outputDir)
		}
		if dryRun {
			return probeWritable(filepath.Dir(outputDir))
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("%w: stat output directory: %w", ErrOutputDirUnavailable, err)
	}
	parent := filepath.Dir(outputDir)
	parentInfo, err := os.Stat(parent)
	if err != nil {
		return fmt.Errorf("%w: output parent directory is not available: %w", ErrOutputDirUnavailable, err)
	}
	if !parentInfo.IsDir() {
		return fmt.Errorf("%w: output parent is not a directory: %s", ErrOutputDirUnavailable, parent)
	}
	return probeWritable(parent)
}

func writeSkillDir(outputDir string, skillMarkdown string, metadataJSON []byte) error {
	parent := filepath.Dir(outputDir)
	tmpDir, err := os.MkdirTemp(parent, "."+safeTempPrefix(filepath.Base(outputDir))+"-")
	if err != nil {
		return fmt.Errorf("create temporary skill directory: %w", err)
	}
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.RemoveAll(tmpDir)
		}
	}()

	if err := os.WriteFile(filepath.Join(tmpDir, "SKILL.md"), []byte(skillMarkdown), 0o644); err != nil {
		return fmt.Errorf("write temporary SKILL.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "loomloom-skill.json"), metadataJSON, 0o644); err != nil {
		return fmt.Errorf("write temporary loomloom-skill.json: %w", err)
	}
	bestEffortSyncDir(tmpDir)
	if err := removeEmptyOutputDirIfPresent(outputDir); err != nil {
		return err
	}
	if err := os.Rename(tmpDir, outputDir); err != nil {
		return fmt.Errorf("publish skill directory: %w", err)
	}
	cleanup = false
	bestEffortSyncDir(parent)
	return nil
}

func removeEmptyOutputDirIfPresent(outputDir string) error {
	info, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%w: stat output directory before publish: %w", ErrOutputDirUnavailable, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrOutputPathNotDir, outputDir)
	}
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("%w: read output directory before publish: %w", ErrOutputDirUnavailable, err)
	}
	if len(entries) > 0 {
		return fmt.Errorf("%w: %s", ErrOutputDirNotEmpty, outputDir)
	}
	if err := os.Remove(outputDir); err != nil {
		return fmt.Errorf("%w: remove empty output directory before publish: %w", ErrOutputDirUnavailable, err)
	}
	return nil
}

func probeWritable(dir string) error {
	tmpDir, err := os.MkdirTemp(dir, ".loomloom-skill-probe-")
	if err != nil {
		return fmt.Errorf("%w: probe write in %s: %w", ErrOutputDirUnavailable, dir, err)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("%w: remove write probe in %s: %w", ErrOutputDirUnavailable, dir, err)
	}
	return nil
}

func bestEffortSyncDir(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	defer f.Close()
	_ = f.Sync()
}

func safeTempPrefix(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "loomloom-skill"
	}
	replacer := strings.NewReplacer(string(filepath.Separator), "-", ":", "-", " ", "-")
	return replacer.Replace(name)
}
