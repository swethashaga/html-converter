package converters

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestConvertHeading tests the conversion of Markdown headings to HTML headings.
func TestConvertHeading(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"# Heading 1", "<h1>Heading 1</h1>"},
		{"## Heading 2", "<h2>Heading 2</h2>"},
		{"###### Heading 6", "<h6>Heading 6</h6>"},
		{"#HeadingWithoutSpace", ""}, // Invalid, should return empty
	}

	for _, c := range cases {
		result, _ := convertHeading(c.input)
		if result != c.expected {
			t.Errorf("convertHeading(%q) == %q, expected %q", c.input, result, c.expected)
		}
	}
}

// TestConvertLink tests the conversion of Markdown links to HTML links.
func TestConvertLink(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"This is a [link](https://example.com)", "This is a <a href=\"https://example.com\">link</a>"},
		{"Just text", "Just text"}, // No link, should remain the same
	}

	for _, c := range cases {
		result, _ := convertLink(c.input)
		if result != c.expected {
			t.Errorf("convertLink(%q) == %q, expected %q", c.input, result, c.expected)
		}
	}
}

// TestConvertEllipsis tests the handling of ellipses.
func TestConvertEllipsis(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"...", "..."},
		{"Not an ellipsis", ""}, // Should return empty since it's not ellipsis
	}

	for _, c := range cases {
		result, _ := convertEllipsis(c.input)
		if result != c.expected {
			t.Errorf("convertEllipsis(%q) == %q, expected %q", c.input, result, c.expected)
		}
	}
}

// TestConvertParagraph tests the conversion of text to HTML paragraphs.
func TestConvertParagraph(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"This is a paragraph.", "<p>This is a paragraph.</p>"},
		// Empty line should return empty
		{"", ""},
		// Ellipsis should not be wrapped in a paragraph
		{"...", ""},
		{"How are you?\nWhat's going on?", "<p>How are you?\nWhat's going on?</p>"},
	}

	for _, c := range cases {
		result, _ := convertParagraph(c.input)
		if result != c.expected {
			t.Errorf("convertParagraph(%q) == %q, expected %q", c.input, result, c.expected)
		}
	}
}

// TestProcessFile tests the processing of a full file.
func TestProcessFile(t *testing.T) {
	// Create a temporary file with some markdown content
	content := `# Title
This is sample markdown for the [Mailchimp](https://www.mailchimp.com) homework assignment.

...

Another paragraph.`

	expectedOutput := `<h1>Title</h1>
<p>This is sample markdown for the <a href="https://www.mailchimp.com">Mailchimp</a> homework assignment.</p>

...

<p>Another paragraph.</p>`

	tmpFile := "test_sample.md"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = ProcessFile(tmpFile)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	if strings.TrimSpace(string(out)) != strings.TrimSpace(expectedOutput) {
		t.Errorf("ProcessFile output was incorrect, got: %s, want: %s", out, expectedOutput)
	}
}
