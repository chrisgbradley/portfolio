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

type Project struct {
	Id          int
	Name        string   `yaml: name`
	Tags        []string `yaml: tags`
	Link        string   `yaml: link`
	Description string   `yaml: description`
}

type Projects = []Project

func newProject(name string, description string, tags []string, link string) Project {
	id++
	return Project{
		Id:          id,
		Name:        name,
		Description: description,
		Tags:        tags,
		Link:        link,
	}
}

func initProjects() Projects {
	entries, err := filepath.Glob(PROJECTS_DATA)
	var projects = Projects{}

	if err != nil {
		panic("dir read fucked")
	}

	for _, entry := range entries {
		p := Project{}
		p.Id = id
		id++
		file, err := os.ReadFile(entry)
		if err != nil {
			panic("file read fucked")
		}
		err = yaml.Unmarshal(file, &p)
		if err != nil {
			panic("yaml parse fucked")
		}
		projects = append(projects, p)
	}

	return projects
}

type Data struct {
	Projects Projects
}

func newData() Data {
	return Data{
		Projects: initProjects(),
	}
}

type Page struct {
	Data Data
}

func newPage() Page {
	return Page{
		Data: newData(),
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	page := newPage()
	e.Renderer = newTemplate()

	e.Static("/css", "css")
	e.Static("/images", "images")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.Logger.Fatal(e.Start(":1324"))
}
