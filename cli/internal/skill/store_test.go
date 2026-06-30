package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteSkillDirPublishesIntoExistingEmptyDirectory(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "skill")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := writeSkillDir(outputDir, "# Test Skill\n", []byte(`{"schema_version":"test"}`+"\n")); err != nil {
		t.Fatalf("writeSkillDir error = %v", err)
	}

	skillBytes, err := os.ReadFile(filepath.Join(outputDir, "SKILL.md"))
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}
	if string(skillBytes) != "# Test Skill\n" {
		t.Fatalf("SKILL.md=%q", string(skillBytes))
	}

	metadataBytes, err := os.ReadFile(filepath.Join(outputDir, "loomloom-skill.json"))
	if err != nil {
		t.Fatalf("read loomloom-skill.json: %v", err)
	}
	if string(metadataBytes) != "{\"schema_version\":\"test\"}\n" {
		t.Fatalf("metadata=%q", string(metadataBytes))
	}
}

func TestWriteSkillDirRejectsExistingNonEmptyDirectory(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "skill")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "existing.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := checkOutputDir(outputDir, false)
	if err == nil {
		t.Fatal("checkOutputDir error = nil, want non-empty directory error")
	}
}
