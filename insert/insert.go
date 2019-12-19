package insert

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func AddLinksToProject(id string) error {
	fn := "./projects/" + id + ".html"

	areLinksInProject, err := CheckIfLinksInPage(fn)
	if err != nil {
		return err
	}
	if areLinksInProject {
		return nil
	}

	styleIndex, err := GetLineNum(fn, "</style>")
	if err != nil {
		return err
	}

	bodyIndex, err := GetLineNum(fn, "</body>")
	if err != nil {
		return err
	}

	err = InsertLine(fn, ".nav-links { position: absolute; top: 15px; left: 15px; }", styleIndex-2)
	if err != nil {
		return err
	}

	err = InsertLine(fn, ".nav-link { font-size: 10pt; color: #bda0a0; margin-left: 5px; text-decoration: none; }", styleIndex-1)
	if err != nil {
		return err
	}

	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='https://www.maxtaylordavi.es'><u>home</u></a><a class='nav-link' href='https://www.maxtaylordavi.es/projects'><u>all projects</u></a></div>", bodyIndex-1)
	if err != nil {
		return err
	}

	return nil
}

func AddLinksToPost(id string) error {
	fn := "./posts/" + id + ".html"

	areLinksInPost, err := CheckIfLinksInPage(fn)
	if err != nil {
		return err
	}
	if areLinksInPost {
		return nil
	}

	styleIndex, err := GetLineNum(fn, "</style>")
	if err != nil {
		return err
	}

	bodyIndex, err := GetLineNum(fn, "</body>")
	if err != nil {
		return err
	}

	err = InsertLine(fn, ".nav-links { position: absolute; top: 15px; left: 15px; }", styleIndex-2)
	if err != nil {
		return err
	}

	err = InsertLine(fn, ".nav-link { font-size: 10pt; color: #bda0a0; margin-left: 5px; text-decoration: none; }", styleIndex-1)
	if err != nil {
		return err
	}

	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='https://www.maxtaylordavi.es'><u>home</u></a><a class='nav-link' href='https://www.maxtaylordavi.es/posts'><u>all posts</u></a></div>", bodyIndex-1)
	if err != nil {
		return err
	}

	return nil
}

func CheckIfLinksInPage(fn string) (bool, error) {
	_, err := GetLineNum(fn, "<div class='nav-links'>")

	if err == nil {
		return true, nil
	}

	if err.Error() == "line not found" {
		return false, nil
	}

	return false, err
}

// get the (0-indexed) number of the line containing the target string
func GetLineNum(fn string, target string) (int, error) {
	f, err := os.Open(fn)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return line, nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// not found
	return 0, errors.New("line not found")
}

func InsertLine(fn, str string, index int) error {
	lines, err := File2lines(fn)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(fn, []byte(fileContent), 0644)
}

func File2lines(fn string) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
