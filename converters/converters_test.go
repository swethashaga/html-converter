package converters

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestProcessFile tests the processing of a full file.
func TestProcessFile(t *testing.T) {
	// Create a test case struct
	type testCase struct {
		input    string
		expected string
	}

	// Define test cases
	testCases := []testCase{
		{
			input: `# Title
This is sample markdown for the [Mailchimp](https://www.mailchimp.com) homework assignment.

...

Another paragraph.`,
			expected: `<h1>Title</h1>
<p>This is sample markdown for the <a href="https://www.mailchimp.com">Mailchimp</a> homework assignment.</p>

...

<p>Another paragraph.</p>`,
		},
		{
			input: `## Subtitle
Just some text with a [link](https://example.com).

...

More text.`,
			expected: `<h2>Subtitle</h2>
<p>Just some text with a <a href="https://example.com">link</a>.</p>

...

<p>More text.</p>`,
		},
		{
			input:    `[Standalone Link](https://example.com)`,
			expected: `<a href="https://example.com">Standalone Link</a>`,
		},
		{
			input:    `...`,
			expected: `...`,
		},
		{
			input: `# Heading
This is another paragraph.`,
			expected: `<h1>Heading</h1>
<p>This is another paragraph.</p>`,
		},
		{
			input: `Empty lines should be preserved.

Another paragraph after empty lines.`,
			expected: `<p>Empty lines should be preserved.</p>

<p>Another paragraph after empty lines.</p>`,
		},
	}

	// Iterate over each test case
	for _, tc := range testCases {
		// Create a temporary file with the input markdown content
		tmpFile := "test_sample.md"
		err := os.WriteFile(tmpFile, []byte(tc.input), 0644)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile) // Ensure the temporary file is removed

		// Capture the output
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Process the file
		err = ProcessFile(tmpFile)
		if err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = old

		// Compare the output with the expected result
		if strings.TrimSpace(string(out)) != strings.TrimSpace(tc.expected) {
			t.Errorf("ProcessFile output was incorrect, got: %s, want: %s", string(out), tc.expected)
		}
	}
}
