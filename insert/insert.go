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

	s0 := ".nav-links { position: fixed; top: 0px; left: 0px; padding: 15px; width: 100%; background-color: white; z-index: 100 } "
	s1 := ".nav-link { display: inline-block; padding-top: 10px; padding-bottom: 10px; width: 100px; font-size: 12pt; color: white; background-color: #9147ff; border-radius: 5px; text-align: center; text-decoration: none; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol; margin-right: 10px; } "
	s2 := ".nav-link:hover { background-image: -webkit-linear-gradient(0deg, #9147ff, #e466bb); -webkit-animation: hue 3s infinite linear; text-decoration: none; } "
	s3 := "@-webkit-keyframes hue { from { -webkit-filter: hue-rotate(0deg); } to { -webkit-filter: hue-rotate(-360deg); } } "
	s4 := ".page-title {display: flex; justify-content: center; width: 100%; margin-top: 100px;} "
	err = InsertLine(fn, s0+s1+s2+s3+s4, styleIndex-1)
	if err != nil {
		return err
	}

	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='/'>home</a><a class='nav-link' href='/projects'>all projects</a></div>", bodyIndex)
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

	s0 := ".nav-links { position: fixed; top: 0px; left: 0px; padding: 15px; width: 100%; background-color: white; z-index: 100 } "
	s1 := ".nav-link { display: inline-block; padding-top: 10px; padding-bottom: 10px; width: 100px; font-size: 12pt; color: white; background-color: #9147ff; border-radius: 5px; text-align: center; text-decoration: none; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol; margin-right: 10px; } "
	s2 := ".nav-link:hover { background-image: -webkit-linear-gradient(0deg, #9147ff, #e466bb); -webkit-animation: hue 3s infinite linear; text-decoration: none; } "
	s3 := "@-webkit-keyframes hue { from { -webkit-filter: hue-rotate(0deg); } to { -webkit-filter: hue-rotate(-360deg); } } "
	s4 := ".page-title {display: flex; justify-content: center; width: 100%; margin-top: 100px;} "
	err = InsertLine(fn, s0+s1+s2+s3+s4, styleIndex-1)
	if err != nil {
		return err
	}

	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='/'>home</a><a class='nav-link' href='/posts'>all posts</a></div>", bodyIndex-1)
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
