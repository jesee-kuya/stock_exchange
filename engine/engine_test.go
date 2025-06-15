package engine

import (
	"os"
	"path/filepath"
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
}
