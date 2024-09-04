package converters

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// convertHeading handles the conversion of Markdown headings (e.g., # Heading) to HTML <h1> to <h6> tags.
func convertHeading(line string) (string, error) {
	for i := 6; i >= 1; i-- {
		prefix := strings.Repeat("#", i) + " "
		if strings.HasPrefix(line, prefix) {
			content := strings.TrimPrefix(line, prefix)
			return fmt.Sprintf("<h%d>%s</h%d>", i, content, i), nil
		}
	}
	return "", nil
}

// convertLink handles the conversion of Markdown links ([text](url)) to HTML <a> tags.
func convertLink(line string) (string, error) {
	linkRegex := regexp.MustCompile(`\[(.+?)\]\((https?:\/\/.+?)\)`)
	if linkRegex.MatchString(line) {
		line = linkRegex.ReplaceAllString(line, `<a href="$2">$1</a>`)
	}
	return line, nil
}

// convertEllipsis handles ellipses (`...`) and returns them as-is.
func convertEllipsis(line string) (string, error) {
	if strings.TrimSpace(line) == "..." {
		return "...", nil
	}
	return "", nil
}

// convertParagraph handles the conversion of any line to HTML <p> tags if it doesn't match any special case (e.g., heading or ellipsis).
func convertParagraph(line string) (string, error) {
	if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") && line != "..." {
		return fmt.Sprintf("<p>%s</p>", line), nil
	}
	return "", nil
}

// ProcessFile reads and converts each line from the input file.
func ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: could not close file %s properly: %v\n", filename, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	var outputLines []string
	var currentParagraph strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check for empty lines to finalize the current paragraph
		if strings.TrimSpace(line) == "" {
			if currentParagraph.Len() > 0 {
				outputLines = append(outputLines, fmt.Sprintf("<p>%s</p>", currentParagraph.String()))
				currentParagraph.Reset()
			}
			// Preserve the empty line
			outputLines = append(outputLines, "")
			continue
		}

		// Process the line using the convertLine function
		converted, err := convertLine(line)
		if err != nil {
			fmt.Printf("error processing line: %v\n", err)
			continue
		}

		// If the line is a block level element (e.g., heading, ellipses), finalize the current paragraph
		if strings.HasPrefix(converted, "<h") || converted == "..." {
			if currentParagraph.Len() > 0 {
				outputLines = append(outputLines, fmt.Sprintf("<p>%s</p>", currentParagraph.String()))
				currentParagraph.Reset()
			}
			outputLines = append(outputLines, converted)
		} else {
			// Accumulate the line into the current paragraph, preserving line breaks
			if currentParagraph.Len() > 0 {
				currentParagraph.WriteString("\n") // Preserve line breaks within paragraphs
			}
			currentParagraph.WriteString(converted)
		}
	}

	// Finalize any remaining paragraph
	if currentParagraph.Len() > 0 {
		outputLines = append(outputLines, fmt.Sprintf("<p>%s</p>", currentParagraph.String()))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Print the converted HTML output with proper formatting
	for _, outputLine := range outputLines {
		fmt.Println(outputLine)
	}

	return nil
}

// convertLine processes each line through all available converters.
func convertLine(line string) (string, error) {
	// Handle empty lines by returning them as-is
	if strings.TrimSpace(line) == "" {
		return "", nil
	}

	// Check for ellipses
	if strings.TrimSpace(line) == "..." {
		return "...", nil
	}

	// Convert headings and check for inline links within the heading
	for i := 6; i >= 1; i-- {
		prefix := strings.Repeat("#", i) + " "
		if strings.HasPrefix(line, prefix) {
			content := strings.TrimPrefix(line, prefix)
			// Convert any inline links within the heading
			content = convertInlineLinks(content)
			return fmt.Sprintf("<h%d>%s</h%d>", i, content, i), nil
		}
	}

	// Convert links in regular paragraphs
	line = convertInlineLinks(line)

	// Return the line as is, without wrapping in <p> tags
	return line, nil
}

// convertInlineLinks handles the conversion of Markdown links ([text](url)) to HTML <a> tags within a line.
func convertInlineLinks(line string) string {
	linkRegex := regexp.MustCompile(`\[(.+?)\]\((https?:\/\/.+?)\)`)
	return linkRegex.ReplaceAllString(line, `<a href="$2">$1</a>`)
}
