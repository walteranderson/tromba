package renderer

import (
	"bytes"
	"html/template"
	"io"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Renderer struct {
	renderer *html.Renderer
}

const MD_EXTENSIONS = parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
const MD_HTML_FLAGS = html.CommonFlags | html.HrefTargetBlank

var renderOpt = html.RendererOptions{Flags: MD_HTML_FLAGS}

func New() *Renderer {
	return &Renderer{
		renderer: html.NewRenderer(renderOpt),
	}
}

func (r *Renderer) Markdown(path string) (string, error) {
	md, err := getFileContent(path)
	if err != nil {
		return "", err
	}
	p := parser.NewWithExtensions(MD_EXTENSIONS)
	doc := p.Parse(md)
	html := markdown.Render(doc, r.renderer)
	return string(html), nil
}

func (r *Renderer) Html(path string) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, "")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
