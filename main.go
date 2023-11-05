package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"maxtaylordavi.es/design"
	"maxtaylordavi.es/insert"
	"maxtaylordavi.es/repository"
)

type Payload = struct {
	Posts        []repository.Post
	Projects     []repository.Project
	Publications []repository.PublicationYear
	Theme        design.Theme
	AllThemes    []design.Theme
	Filtered     bool
	Title        string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("home.gohtml", "projects.gohtml", "posts.gohtml", "publications.gohtml"))
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func idxToLetter(i int) string {
	return string('A' + i)
}

func oddOrEven(i int) string {
	if i%2 == 0 {
		return "even"
	}
	return "odd"
}

var fm = template.FuncMap{"fdate": formatDate, "i2l": idxToLetter, "oddOrEven": oddOrEven}

func serveImage(path string, w http.ResponseWriter) {
	img, err := os.Open(path)
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer img.Close()
	io.Copy(w, img)
}

func registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}
		json.NewEncoder(w).Encode(payload)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		tag := r.URL.Query().Get("tag")
		theme := design.GetTheme(r.URL.Query().Get("theme"))

		postr := repository.PostRepository{}

		var posts []repository.Post
		var title string
		var err error
		if tag == "" {
			posts, err = postr.All()
			title = "All posts"
		} else {
			posts, err = postr.WithTag(tag)
			title = fmt.Sprintf("Posts tagged <a href='/posts?tag=%s&theme=%s' class='tag %s title'>%s</a>", tag, theme.Name, tag, tag)
		}

		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the posts page
		data := Payload{
			Posts:    posts,
			Theme:    theme,
			Title:    title,
			Filtered: tag != "",
		}

		_ = tpl.ExecuteTemplate(w, "posts.gohtml", data)
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		tag := r.URL.Query().Get("tag")
		theme := design.GetTheme(r.URL.Query().Get("theme"))

		projr := repository.ProjectRepository{}

		var projects []repository.Project
		var title string
		var err error
		if tag == "" {
			projects, err = projr.All()
			title = "All projects"
		} else {
			projects, err = projr.WithTag(tag)
			title = fmt.Sprintf("Projects tagged <a href='/projects?tag=%s&theme=%s' class='tag %s title'>%s</a>", tag, theme.Name, tag, tag)
		}

		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the projects page
		data := Payload{
			Projects: projects,
			Theme:    theme,
			Title:    title,
			Filtered: tag != "",
		}

		_ = tpl.ExecuteTemplate(w, "projects.gohtml", data)
	})

	mux.HandleFunc("/publications", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		tag := r.URL.Query().Get("tag")
		theme := design.GetTheme(r.URL.Query().Get("theme"))

		pubr := repository.PublicationRepository{}

		var publications []repository.PublicationYear
		var title string
		var err error
		if tag == "" {
			publications, err = pubr.All()
			title = "All publications"
		} else {
			publications, err = pubr.WithTag(tag)
			title = fmt.Sprintf("Publications tagged <a href='/publications?tag=%s&theme=%s' class='tag %s title'>%s</a>", tag, theme.Name, tag, tag)
		}

		if err != nil {
			log.Fatalln("error getting publications: ", err)
		}

		// serve the publications page
		data := Payload{
			Publications: publications,
			Theme:        theme,
			Title:        title,
			Filtered:     tag != "",
		}

		err = tpl.ExecuteTemplate(w, "publications.gohtml", data)
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		id := r.URL.Query().Get("id")
		themeName := r.URL.Query().Get("theme")

		b, err := insert.ReadFileAndInjectStuff("./posts/"+id+"/main.html", themeName)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})

	mux.HandleFunc("/cv", func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("./CV.pdf")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(f))
	})

	mux.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/media/")
		parts := strings.Split(path, "/")

		ct := "image/png"
		if strings.Contains(parts[1], ".svg") {
			ct = "image/svg+xml"
		}
		w.Header().Set("Content-Type", ct)

		f, err := os.Open("./posts/" + parts[0] + "/media/" + parts[1])
		if err != nil {
			log.Fatal(err) // perhaps handle this nicer
		}

		defer f.Close()
		io.Copy(w, f)
	})

	mux.HandleFunc("/styles/", func(w http.ResponseWriter, r *http.Request) {
		css, err := os.Open("." + strings.TrimSuffix(r.URL.Path, "/"))
		if err != nil {
			log.Fatal(err) // perhaps handle this nicer
		}
		defer css.Close()
		w.Header().Set("Content-Type", "text/css")
		io.Copy(w, css)
	})

	mux.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		serveImage("."+strings.TrimSuffix(r.URL.Path, "/"), w)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		serveImage("./favicon.ico", w)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		allThemes := design.GetAllThemes()
		theme := design.GetTheme(r.URL.Query().Get("theme"))

		data := Payload{
			Theme:     theme,
			AllThemes: allThemes,
		}

		// serve the homepage
		_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
	})

	return mux
}

func main() {
	server := http.Server{
		Addr:         ":8000",
		Handler:      registerRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()
}
