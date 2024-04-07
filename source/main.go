package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"
)

const PROJECTS_DATA = "./data/projects/*.yaml"

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id = 0

// TODO: Break Project and Projects types, functions, etc. into separate file
// FIXME: Fix the Projects/Project nomensclature for clarifcation -- impossible to read at the moment
type Project struct {
	Id          int
	Name        string   `yaml: name`
	Tags        []string `yaml: tags`
	Link        string   `yaml: link`
	Description string   `yaml: description`
}

type Projects struct {
	Projects []Project
}

func initProjects() Projects {
	entries, err := filepath.Glob(PROJECTS_DATA)
	var projects = Projects{Projects: []Project{}}

	if err != nil {
		panic("Directory is empty or other error found.")
	}

	for _, entry := range entries {
		p := Project{}
		p.Id = id
		id++
		file, err := os.ReadFile(entry)
		if err != nil {
			panic("File unable to be read.")
		}
		err = yaml.Unmarshal(file, &p)
		if err != nil {
			panic("Unable to unmarshal/parse the yaml file.")
		}
		projects.Projects = append(projects.Projects, p)
	}

	return projects
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	projects := initProjects()
	e.Renderer = newTemplate()

	e.Static("/", "www")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", projects)
	})

	e.Logger.Fatal(e.Start(":1324"))
}
