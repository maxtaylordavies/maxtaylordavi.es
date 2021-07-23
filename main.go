package main

import (
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
	Posts     []repository.Post
	Projects  []repository.Project
	Theme     design.Theme
	AllThemes []design.Theme
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("home.gohtml", "projects.gohtml", "posts.gohtml"))
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func idxToLetter(i int) string {
	return string('A' + i)
}

var fm = template.FuncMap{"fdate": formatDate, "i2l": idxToLetter}

func registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		theme := design.GetTheme(r.URL.Query().Get("theme"))

		postr := repository.PostRepository{}
		posts, err := postr.All()

		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the posts page
		data := Payload{
			Posts: posts,
			Theme: theme,
		}

		_ = tpl.ExecuteTemplate(w, "posts.gohtml", data)
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		theme := design.GetTheme(r.URL.Query().Get("theme"))

		projr := repository.ProjectRepository{}
		projects, err := projr.All()
		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the projects page
		data := Payload{
			Projects: projects,
			Theme:    theme,
		}

		_ = tpl.ExecuteTemplate(w, "projects.gohtml", data)
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
		img, err := os.Open("." + strings.TrimSuffix(r.URL.Path, "/"))
		if err != nil {
			log.Fatal(err) // perhaps handle this nicer
		}
		defer img.Close()
		// w.Header().Set("Content-Type", "text/css")
		io.Copy(w, img)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		allThemes := design.GetAllThemes()
		theme := design.GetTheme(r.URL.Query().Get("theme"))

		postr := repository.PostRepository{}
		projr := repository.ProjectRepository{}

		recentPosts, err := postr.Recent()
		if err != nil {
			log.Fatalln("error getting recent posts: ", err)
		}

		recentProjects, err := projr.Recent()
		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		data := Payload{
			Posts:     recentPosts,
			Projects:  recentProjects,
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
