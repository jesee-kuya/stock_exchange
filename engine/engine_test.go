package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLog(t *testing.T) {
	t.Run("file format", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "simulation.log")

		// mock schedule data
		engine := &Engine{
			Schedule: []string{
				"1:process_a",
				"1:process_b",
				"2:process_c",
				"3:process_a",
			},
		}
		err := engine.SaveLog(testPath)
		if err != nil {
			t.Fatalf("SaveLog returned an error: %v", err)
		}

		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		want := "1:process_a\n1:process_b\n2:process_c\n3:process_a"
		if string(content) != want {
			t.Errorf("File content mismatch. got %v but expected %v", content, want)
		}
	})

	t.Run("empty path string", func(t *testing.T) {
		engine := &Engine{
			Schedule: []string{"1:test_process"},
		}

		err := engine.SaveLog("")

		if err == nil {
			t.Error("SaveLog() should return error for empty path")
		}
	})

	t.Run("nil schedule", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "nil.log")

		engine := &Engine{
			Schedule: nil, // Nil schedule
		}

		err := engine.SaveLog(testPath)
		if err != nil {
			t.Errorf("SaveLog() should handle nil schedule gracefully: %v", err)
		}

		// Verify empty file was created
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if len(strings.TrimSpace(string(content))) != 0 {
			t.Error("Expected empty file for nil schedule")
		}
	})
}
