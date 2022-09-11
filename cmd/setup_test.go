package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestWriteFile(t *testing.T) {
	testData := []byte("test")

	err := writeFile("test.txt", testData)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	// Read the file and compare the contents.
	f, err := os.Open("test.txt")
	if err != nil {
		t.Errorf("Error opening file: %v", err)
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	if !bytes.Contains(b, testData) {
		t.Errorf("Expected %v, got %v", string(testData), string(b))
	}

	// Remove the file.
	err = os.Remove("test.txt")
	if err != nil {
		t.Errorf("Error removing file: %v", err)
	}

	// Try to read the file again.
	f, err = os.Open("test.txt")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}
