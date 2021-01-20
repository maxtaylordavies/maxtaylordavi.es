package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/maxtaylordavies/maxtaylordavi.es/insert"
	"github.com/maxtaylordavies/maxtaylordavi.es/repository"
	"golang.org/x/crypto/acme/autocert"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("home.gohtml", "projects.gohtml", "posts.gohtml"))
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

var fm = template.FuncMap{"fdate": formatDate}

func registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		theme := r.URL.Query().Get("theme")

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

		data := struct {
			Posts    []repository.Post
			Projects []repository.Project
			Theme    string
		}{
			Posts:    recentPosts,
			Projects: recentProjects,
			Theme:    theme,
		}

		// serve the homepage
		_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		postr := repository.PostRepository{}
		posts, err := postr.All()

		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the posts page
		_ = tpl.ExecuteTemplate(w, "posts.gohtml", posts)
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		projr := repository.ProjectRepository{}
		projects, err := projr.All()
		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the projects page
		_ = tpl.ExecuteTemplate(w, "projects.gohtml", projects)
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		id := r.URL.Query().Get("id")
		if id == "" {
			log.Println("no id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := insert.AddLinksToPost(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = insert.AddTagsToPostOrProject(id, "post")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f, err := ioutil.ReadFile("./posts/" + id + ".html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(f))
	})

	mux.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := insert.AddLinksToProject(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = insert.AddTagsToPostOrProject(id, "project")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f, err := ioutil.ReadFile("./projects/" + id + ".html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(f))
	})

	mux.HandleFunc("/cv", func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("./CV.pdf")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(f))
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

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
	}

	server := http.Server{
		Addr: ":https",
		// Addr: ":80",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler:      registerRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
