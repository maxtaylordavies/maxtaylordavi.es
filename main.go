package main

import (
	"./repository"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("home.gohtml", "projects.gohtml", "posts.gohtml", "project.gohtml", "post.gohtml"))
}

func formatDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func getIntro(body string) string {
	slc := strings.SplitAfter(body, ". ")[0:3]
	str := strings.Join(slc, "")
	str = str[0:len(str)-1] + ".."
	return str
}

var fm = template.FuncMap{"fdate": formatDate, "intro": getIntro}

func registerRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		postr := repository.PostRepository{db}
		projr := repository.ProjectRepository{db}
		recentPosts, err := postr.Recent()
		if err != nil {
			log.Fatalln("error getting recent posts: ", err)
		}
		recentProjects, err := projr.Recent()
		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		data := struct{
			Posts []repository.Post
			Projects []repository.Project
		}{
			Posts: recentPosts,
			Projects: recentProjects,
		}

		// serve the homepage
		_ = tpl.ExecuteTemplate(w, "home.gohtml", data)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		postr := repository.PostRepository{db}
		posts, err := postr.All()
		if err != nil {
			log.Fatalln("error getting recent projects: ", err)
		}

		// serve the posts page
		_ = tpl.ExecuteTemplate(w, "posts.gohtml", posts)
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		projr := repository.ProjectRepository{db}
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

		postr := repository.PostRepository{db}
		post, err := postr.One(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = tpl.ExecuteTemplate(w, "post.gohtml", post)
	})

	mux.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		projr := repository.ProjectRepository{db}
		project, err := projr.One(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = tpl.ExecuteTemplate(w, "project.gohtml", project)
	})

	return mux
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	mux := registerRoutes(db)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
