package project

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/walteranderson/tromba/internal/config"
)

type Project struct {
	Config *config.Config
	Pages  []*Page
	wg     sync.WaitGroup
}

type Page struct {
	Path        string // absolute path to file
	Url         string // /blog/article-one, /nested/one/two
	Filename    string // +page.html, +article-one.md
	Ext         string // html, md
	HtmlContent string
}

func Build(c *config.Config) (*Project, error) {
	proj := &Project{
		Config: c,
		Pages:  []*Page{},
	}

	err := proj.walkPages()
	if err != nil {
		return nil, err
	}

	for _, page := range proj.Pages {
		proj.wg.Add(1)
		go proj.processPage(page)
	}
	proj.wg.Wait()

	outDir := filepath.Join(proj.Config.ProjectPath, proj.Config.ProjectConfig.Output)
	for _, page := range proj.Pages {
		outPath := filepath.Join(outDir, page.Url)
		_, err := os.Stat(outPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(outPath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}

		file, err := os.Create(filepath.Join(outPath, "index.html"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(page.HtmlContent)
		if err != nil {
			log.Fatal(err)
		}
	}

	return proj, nil
}

// TODO: how to handle errors? channel?
func (p *Project) processPage(page *Page) {
	defer p.wg.Done()
	switch page.Ext {
	case "md":
		err := page.renderMarkdown()
		if err != nil {
			// TODO
			log.Println(err)
		}
	case "html":
		err := page.renderHtml()
		if err != nil {
			// TODO
			log.Println(err)
		}
	default:
		log.Printf("Unsupported file extension: %s\n", page.Ext)
	}
}

func (p *Project) walkPages() error {
	rootPath := filepath.Join(p.Config.ProjectPath, p.Config.ProjectConfig.Input)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pagePath := strings.Split(path, p.Config.ProjectPath)[1]
		pagePath = strings.TrimPrefix(pagePath, "/"+p.Config.ProjectConfig.Input)

		ext := filepath.Ext(pagePath)
		if ext == "" {
			return nil
		}
		ext = strings.TrimPrefix(ext, ".")

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

		page := &Page{
			Path:     path,
			Url:      pagePath,
			Filename: filename,
			Ext:      ext,
		}
		p.Pages = append(p.Pages, page)
		return nil
	})
	return err
}

func (page *Page) renderMarkdown() error {
	md, err := getFileContent(page.Path)
	if err != nil {
		return err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	html := markdown.Render(doc, renderer)
	page.HtmlContent = string(html)
	return nil
}

func (page *Page) renderHtml() error {
	if page.Ext != "html" {
		return fmt.Errorf("Trying to render HTML on the wrong filetype, got %s", page.Ext)
	}
	tmpl, err := template.ParseFiles(page.Path)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, "")
	if err != nil {
		return err
	}
	page.HtmlContent = buf.String()
	return nil
}

func getFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
