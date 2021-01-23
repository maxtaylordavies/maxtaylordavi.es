package insert

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/maxtaylordavies/maxtaylordavi.es/design"
)

func ReadFileAndInjectStuff(fn, themeName string) ([]byte, error) {
	// read html file into a string
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return b, err
	}
	s := string(b)

	// add tags
	s = AddTags(s)

	// add nav links
	s = AddNavLinks(s, strings.Contains(fn, "project"), themeName)

	return []byte(s), nil
}

func AddTags(s string) string {
	i := strings.Index(s, "<em>")
	j := strings.Index(s, "</em>")
	rawTags := strings.Split(s[i+15:j], " ")
	styledTags := fmt.Sprintf("<div class='tags'> <em>%s</em>", s[i+4:i+14])
	for _, tag := range rawTags {
		styledTags += fmt.Sprintf("<div class='tag %s'>%s</div>", tag, tag)
	}
	styledTags += "</div>"
	s = s[:i] + styledTags + s[j+5:]

	i = strings.Index(s, "</head>")
	s = s[:i] + fmt.Sprintf("<link rel='stylesheet' href='/styles/tags.css'/><meta name='tags' content='%s'/>", strings.Join(rawTags, " ")) + s[i:]

	return s
}

func AddNavLinks(s string, project bool, themeName string) string {
	backLink := fmt.Sprintf("/posts?theme=%s", themeName)
	if project {
		backLink = fmt.Sprintf("/projects?theme=%s", themeName)
	}

	backLinkText := "all posts"
	if project {
		backLinkText = "all projects"
	}

	// add elements
	i := strings.Index(s, "</body>")
	s = s[:i] + fmt.Sprintf("<div class='nav-links'><a class='nav-link' href='/?theme=%s'>home</a><a class='nav-link' href=%s>%s</a></div>", themeName, backLink, backLinkText) + s[i:]

	// add css
	i = strings.Index(s, "</head>")
	s = s[:i] + "<link rel='stylesheet' href='/styles/links.css'/>" + s[i:]

	// add theme-specific styles
	theme := design.GetTheme(themeName)
	i = strings.Index(s, "</style>")
	s = s[:i] + fmt.Sprintf("html {background: %s} .nav-link {background: %s}", theme.Background, theme.Color) + s[i:]

	return s
}

// func AddLinksToProject(id string) error {
// 	fn := "./projects/" + id + ".html"

// 	areLinksInProject, err := CheckIfLinksInPage(fn)
// 	if err != nil {
// 		return err
// 	}
// 	if areLinksInProject {
// 		return nil
// 	}

// 	styleIndex, err := GetLineNum(fn, "</style>")
// 	if err != nil {
// 		return err
// 	}
// 	bodyIndex, err := GetLineNum(fn, "</body>")
// 	if err != nil {
// 		return err
// 	}

// 	s0 := ".nav-links { position: fixed; top: 0px; left: 0px; padding: 15px; width: 100%; z-index: 100 } "
// 	s1 := ".nav-link { display: inline-block; padding-top: 10px; padding-bottom: 10px; width: 100px; font-size: 12pt; color: white; background-color: #9147ff; border-radius: 5px; text-align: center; text-decoration: none; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol; margin-right: 10px; } "
// 	s2 := ".nav-link:hover { background-image: -webkit-linear-gradient(0deg, #9147ff, #e466bb); -webkit-animation: hue 3s infinite linear; text-decoration: none; } "
// 	s3 := "@-webkit-keyframes hue { from { -webkit-filter: hue-rotate(0deg); } to { -webkit-filter: hue-rotate(-360deg); } } "
// 	s4 := ".page-title {display: flex; justify-content: center; width: 100%; margin-top: 100px;} "
// 	err = InsertLine(fn, s0+s1+s2+s3+s4, styleIndex-1)
// 	if err != nil {
// 		return err
// 	}

// 	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='/'>home</a><a class='nav-link' href='/projects'>all projects</a></div>", bodyIndex)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func AddLinksToPost(id string) error {
// 	fn := "./posts/" + id + ".html"

// 	areLinksInPost, err := CheckIfLinksInPage(fn)
// 	if err != nil {
// 		return err
// 	}
// 	if areLinksInPost {
// 		return nil
// 	}

// 	styleIndex, err := GetLineNum(fn, "</style>")
// 	if err != nil {
// 		return err
// 	}
// 	bodyIndex, err := GetLineNum(fn, "</body>")
// 	if err != nil {
// 		return err
// 	}

// 	s0 := ".nav-links { position: fixed; top: 0px; left: 0px; padding: 15px; width: 100%; background-color: white; z-index: 100 } "
// 	s1 := ".nav-link { display: inline-block; padding-top: 10px; padding-bottom: 10px; width: 100px; font-size: 12pt; color: white; background-color: #9147ff; border-radius: 5px; text-align: center; text-decoration: none; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol; margin-right: 10px; } "
// 	s2 := ".nav-link:hover { background-image: -webkit-linear-gradient(0deg, #9147ff, #e466bb); -webkit-animation: hue 3s infinite linear; text-decoration: none; } "
// 	s3 := "@-webkit-keyframes hue { from { -webkit-filter: hue-rotate(0deg); } to { -webkit-filter: hue-rotate(-360deg); } } "
// 	s4 := ".page-title {display: flex; justify-content: center; width: 100%; margin-top: 100px;} "
// 	err = InsertLine(fn, s0+s1+s2+s3+s4, styleIndex-1)
// 	if err != nil {
// 		return err
// 	}

// 	err = InsertLine(fn, "<div class='nav-links'><a class='nav-link' href='/'>home</a><a class='nav-link' href='/posts'>all posts</a></div>", bodyIndex-1)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // get the (0-indexed) number of the line containing the target string
// func GetLineNum(fn string, target string) (int, error) {
// 	f, err := os.Open(fn)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()

// 	// Splits on newlines by default.
// 	scanner := bufio.NewScanner(f)
// 	line := 0
// 	for scanner.Scan() {
// 		if strings.Contains(scanner.Text(), target) {
// 			return line, nil
// 		}
// 		line++
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return 0, err
// 	}

// 	// not found
// 	return 0, errors.New("line not found")
// }

// func InsertLine(fn, str string, index int) error {
// 	lines, err := File2lines(fn)
// 	if err != nil {
// 		return err
// 	}

// 	fileContent := ""
// 	for i, line := range lines {
// 		if i == index {
// 			fileContent += str
// 		}
// 		fileContent += line
// 		fileContent += "\n"
// 	}

// 	return ioutil.WriteFile(fn, []byte(fileContent), 0644)
// }

// func File2lines(fn string) ([]string, error) {
// 	f, err := os.Open(fn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	return LinesFromReader(f)
// }

// func LinesFromReader(r io.Reader) ([]string, error) {
// 	var lines []string
// 	scanner := bufio.NewScanner(r)
// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}

// 	return lines, nil
// }
