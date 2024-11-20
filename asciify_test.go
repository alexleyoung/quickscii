package quickscii

import (
	"os"
	"strings"
	"testing"
)

func TestAsciify(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		width       int
		height      int
		charset     string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty path",
			path:        "",
			width:       10,
			height:      10,
			charset:     "mix",
			wantErr:     true,
			errContains: "invalid path",
		},
		{
			name:        "invalid dimensions",
			path:        "testdata/test.png",
			width:       -1,
			height:      10,
			charset:     "mix",
			wantErr:     true,
			errContains: "invalid dimensions",
		},
		{
			name:        "invalid charset",
			path:        "testdata/test.png",
			width:       10,
			height:      10,
			charset:     "invalid",
			wantErr:     true,
			errContains: "invalid charset",
		},
		{
			name:        "nonexistent file",
			path:        "nonexistent.png",
			width:       10,
			height:      10,
			charset:     "mix",
			wantErr:     true,
			errContains: "failed to read image",
		},
	}

	// Create test directory and sample image
	err := os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Asciify(tt.path, tt.width, tt.height, tt.charset)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func TestAsciifyToImage(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		width       int
		height      int
		charset     string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty input path",
			inputPath:   "",
			outputPath:  "testdata/output.png",
			width:       10,
			height:      10,
			charset:     "mix",
			wantErr:     true,
			errContains: "invalid path",
		},
		{
			name:        "invalid dimensions",
			inputPath:   "testdata/test.png",
			outputPath:  "testdata/output.png",
			width:       -1,
			height:      10,
			charset:     "mix",
			wantErr:     true,
			errContains: "invalid dimensions",
		},
		{
			name:        "invalid charset",
			inputPath:   "testdata/test.png",
			outputPath:  "testdata/output.png",
			width:       10,
			height:      10,
			charset:     "invalid",
			wantErr:     true,
			errContains: "invalid charset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AsciifyToImage(tt.inputPath, tt.outputPath, tt.width, tt.height, tt.charset)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check if output file exists
			if _, err := os.Stat(tt.outputPath); os.IsNotExist(err) {
				t.Error("Output file was not created")
			}
		})
	}
}

// TestCharsets tests ASCII art generation with different character sets
func TestCharsets(t *testing.T) {
	charsets := []string{"block", "poly", "mix"}
	
	// Create a test image first
	testImagePath := "testdata/test.png"
	
	for _, charset := range charsets {
		t.Run(charset, func(t *testing.T) {
			result, err := Asciify(testImagePath, 10, 10, charset)
			if err != nil {
				t.Errorf("Failed to generate ASCII art with charset %s: %v", charset, err)
				return
			}
			
			if result == "" {
				t.Errorf("Empty result for charset %s", charset)
			}
		})
	}
}
