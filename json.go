package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Makes the "updated" field up-to-date.
// Can't use json marshaling because it introduces changes to the
// file structure.
func editUpdatedField(fname, to string) error {
	var updatedLine = `  "updated": `
	lines, err := readLines(fname)
	if err != nil {
		return err
	}
	for idx, line := range lines {
		if strings.Contains(line, updatedLine) {
			updated := updatedLine + to + ","
			lines[idx] = updated
		}
	}

	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	w.Flush()
	return nil
}

// Reads fname and inserts "content" before a given line
func insertIntoJsonFile(fname, content string, before int) error {
	lines, err := readLines(fname)
	if err != nil {
		panic(err)
	}

	var linesPre []string
	var linesPost []string

	for idx, line := range lines {
		if idx < before {
			linesPre = append(linesPre, line)
		} else {
			linesPost = append(linesPost, line)
		}
	}

	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 1. Write lines before our content
	w := bufio.NewWriter(file)
	for idx, line := range linesPre {
		if idx == len(linesPre)-1 { // last line
			if !strings.HasSuffix(line, ",") {
				line += ","
			}
		}
		fmt.Fprintln(w, line)
	}

	// 2. Write new content
	fmt.Fprint(w, content)

	// 3. Write post content
	for idx, line := range linesPost {
		// check if a , is needed here
		if idx == 0 {
			if strings.Contains(line, "]") {
				// this is the final build
				fmt.Fprintln(w, "")
			} else {
				fmt.Fprintln(w, ",")
			}
		}
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}

func scanFileForSubstring(fname, substring string) int {
	lines, err := readLines(fname)
	if err != nil {
		panic(err)
	}

	for idx, line := range lines {
		if strings.Contains(line, substring) {
			return idx
		}
	}

	return -1
}

func searchFuncForVersion(version semver.Version) func(string) bool {
	var myFunc = func(line string) bool {
		if strings.Contains(line, `      "name": "`) {

			split := strings.Split(line, ":")
			clean := split[1][2 : len(split[1])-1]
			v2, _ := semver.Make(clean)

			return version.EQ(v2)
		}
		return false
	}

	return myFunc
}

// Scans a file by passing each line to the function passed as second argument
// Returns the position (in lines) in the file where it was found.
func scanFileFor(fname string, searchFunc func(string) bool) int {
	lines, err := readLines(fname)
	if err != nil {
		panic(err)
	}

	for idx, line := range lines {
		if searchFunc(line) {
			return idx
		}
	}

	return -1
}
