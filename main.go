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

	"github.com/maxtaylordavies/maxtaylordavi.es/design"
	"github.com/maxtaylordavies/maxtaylordavi.es/insert"
	"github.com/maxtaylordavies/maxtaylordavi.es/repository"
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

		b, err := insert.ReadFileAndInjectStuff("./posts/"+id+".html", themeName)

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

	mux.HandleFunc("/thesis", func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("./thesis/thesis.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(f))
	})

	mux.HandleFunc("/figures/", func(w http.ResponseWriter, r *http.Request) {
		img, err := os.Open("./thesis" + strings.TrimSuffix(r.URL.Path, "/"))
		if err != nil {
			log.Fatal(err) // perhaps handle this nicer
		}
		defer img.Close()
		io.Copy(w, img)
	})

	mux.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		img, err := os.Open("." + strings.TrimSuffix(r.URL.Path, "/"))
		if err != nil {
			log.Fatal(err) // perhaps handle this nicer
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, img)
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

	return mux
}

func main() {

	// certManager := autocert.Manager{
	// 	Prompt: autocert.AcceptTOS,
	// 	Cache:  autocert.DirCache("certs"),
	// }

	server := http.Server{
		// Addr: ":https",
		Addr: ":8000",
		// TLSConfig: &tls.Config{
		// 	GetCertificate: certManager.GetCertificate,
		// },
		Handler:      registerRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()

	// go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	// server.ListenAndServeTLS("", "")
}
