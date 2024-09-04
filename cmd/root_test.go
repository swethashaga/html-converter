package cmd

import (
	"html-converter/converters"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper function to reset rootCmd between tests
func resetRootCmd() {
	rootCmd = &cobra.Command{
		Use:   "htmlConverter <filename>",
		Short: "Markdown to HTML converter",
		Long:  `A command-line tool to convert Markdown format to HTML.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := converters.ProcessFile(args[0]); err != nil {
				cmd.PrintErrf("Error processing file: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// TestRootCmd tests the root command execution.
func TestRootCmd(t *testing.T) {
	resetRootCmd()

	content := `# Title
This is sample markdown for the [Mailchimp](https://www.mailchimp.com) homework assignment.

...

Another paragraph.`

	expectedOutput := `<h1>Title</h1>
<p>This is sample markdown for the <a href="https://www.mailchimp.com">Mailchimp</a> homework assignment.</p>

...

<p>Another paragraph.</p>`

	tmpFile, err := os.CreateTemp("", "test_sample_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	rootCmd.SetArgs([]string{tmpFile.Name()})

	// Capture the output using os.Pipe
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute root command: %v", err)
	}

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	output := string(out)

	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
