package converters

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

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
		line := strings.TrimSpace(scanner.Text())

		// Check for empty lines to finalize the current paragraph
		if line == "" {
			if currentParagraph.Len() > 0 {
				// Finalize and append the current paragraph
				outputLines = append(outputLines, fmt.Sprintf("<p>%s</p>", currentParagraph.String()))
				currentParagraph.Reset()
			}
			outputLines = append(outputLines, "") // Preserve the empty line
			continue
		}

		// Process the line using the convertLine function
		converted, err := convertLine(line)
		if err != nil {
			fmt.Printf("error processing line: %v\n", err)
			continue
		}

		// Handle standalone links separately
		if isStandaloneLink(line) {
			if currentParagraph.Len() > 0 {
				outputLines = append(outputLines, fmt.Sprintf("<p>%s</p>", currentParagraph.String()))
				currentParagraph.Reset()
			}
			outputLines = append(outputLines, converted)
			continue
		}

		// If the line is a block-level element (e.g., heading, ellipses), finalize the current paragraph
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

	// If it's a standalone link, return the link without wrapping it in <p> tags
	if isStandaloneLink(line) {
		return convertInlineLinks(line), nil
	}

	// Convert inline links for regular text
	return convertInlineLinks(line), nil
}

// isStandaloneLink checks if the line is a standalone link (i.e., a link occupying the whole line).
func isStandaloneLink(line string) bool {
	linkRegex := regexp.MustCompile(`^\[.+?\]\(https?:\/\/.+?\)$`)
	return linkRegex.MatchString(line)
}

// convertInlineLinks handles the conversion of Markdown links ([text](url)) to HTML <a> tags within a line.
func convertInlineLinks(line string) string {
	linkRegex := regexp.MustCompile(`\[(.+?)\]\((https?:\/\/.+?)\)`)
	return linkRegex.ReplaceAllString(line, `<a href="$2">$1</a>`)
}
