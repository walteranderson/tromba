package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/walteranderson/tromba/internal/config"
)

type Project struct {
	Config *config.Config
	Pages  []Page
}

type Page struct {
	Path     string // absolute path to file
	Url      string // /blog/article-one, /nested/one/two
	Filename string // +page.html, +article-one.md
}

func Build(c *config.Config) (*Project, error) {
	proj := &Project{
		Config: c,
	}

	rootPath := filepath.Join(c.ProjectPath, c.ProjectConfig.Input)
	var pages []Page
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pagePath := strings.Split(path, c.ProjectPath)[1]
		pagePath = strings.TrimPrefix(pagePath, "/"+c.ProjectConfig.Input)

		ext := filepath.Ext(pagePath)
		if ext == "" {
			return nil
		}
		_, filename := filepath.Split(pagePath)
		if strings.HasPrefix(filename, "+") {
			// /blog/+page.html => /blog
			pagePath = strings.Split(pagePath, "+")[0]
		} else {
			// /blog/article-one.md => /blog/article-one
			pagePath = strings.Split(pagePath, ".")[0]
		}

		// trim trailing slash
		if pagePath != "/" && strings.HasSuffix(pagePath, "/") {
			pagePath = strings.TrimSuffix(pagePath, "/")
		}

		p := Page{
			Path:     path,
			Url:      pagePath,
			Filename: filename,
		}
		pages = append(pages, p)
		return nil
	})

	proj.Pages = pages

	if err != nil {
		return nil, err
	}

	fmt.Printf("proj: %v\n", proj)
	return proj, nil
}
