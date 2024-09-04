# HTML Converter Tool

This is a command-line tool written in Go that converts Markdown files to HTML. It supports basic Markdown elements such as headings, links, ellipses, and paragraphs. The tool reads a Markdown file as input and outputs the corresponding HTML.

### Steps to Install

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/swethashaga/html-converter.git
    ```

2. Navigate to the project directory:

    ```bash
    cd html-converter
    ```

3. Build the tool using Go:

    ```bash
    go build -o htmlConverter
    ```

## Usage

Once the tool is built, you can use it from the command line to convert Markdown files to HTML.

### Basic Usage

To convert a Markdown file to HTML, run the following command:

```bash
./htmlConverter path/to/your-markdown-file
