package generator

import (
    "bytes"
    "fmt"
    "html/template"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/ast"
    meta "github.com/yuin/goldmark-meta"
    "github.com/yuin/goldmark/parser"
    "github.com/yuin/goldmark/text"
)

type Page struct {
    Title     string
    Subtitle  string
    Slug      string
    TOC       template.HTML
    Content   template.HTML
    CreatedAt string
    UpdatedAt string
}

type Config struct {
    ContentDir  string
    TemplateDir string
    StaticDir   string
    OutputDir   string
    BaseURL     string
}

var requiredPages = []struct {
    file    string
    outPath string
}{
    {"home.md", "index.html"},
    {"privacy.md", filepath.Join("privacy", "index.html")},
    {"terms.md", filepath.Join("terms", "index.html")},
}

func Run(cfg Config) error {
    tmpl, err := template.ParseFiles(
        filepath.Join(cfg.TemplateDir, "layout.html"),
        filepath.Join(cfg.TemplateDir, "page.html"),
    )
    if err != nil {
        return fmt.Errorf("parse templates: %w", err)
    }

    md := goldmark.New(
        goldmark.WithExtensions(meta.Meta),
        goldmark.WithParserOptions(
            parser.WithAutoHeadingID(),
        ),
    )

    for _, rp := range requiredPages {
        srcPath := filepath.Join(cfg.ContentDir, "pages", rp.file)
        src, err := os.ReadFile(srcPath)
        if err != nil {
            return fmt.Errorf("missing required page %q: %w", srcPath, err)
        }

        page, err := parsePage(md, src)
        if err != nil {
            return fmt.Errorf("parse %q: %w", rp.file, err)
        }

        page.CreatedAt, page.UpdatedAt = gitDates(srcPath)

        outPath := filepath.Join(cfg.OutputDir, rp.outPath)
        if err := writeHTML(tmpl, page, outPath); err != nil {
            return fmt.Errorf("write %q: %w", outPath, err)
        }

        fmt.Printf("  wrote %s\n", outPath)
    }

    if err := copyStatic(cfg.StaticDir, cfg.OutputDir); err != nil {
        return fmt.Errorf("copy static: %w", err)
    }

    return nil
}

func parsePage(md goldmark.Markdown, src []byte) (Page, error) {
    ctx := parser.NewContext()
    reader := text.NewReader(src)
    doc := md.Parser().Parse(reader, parser.WithContext(ctx))

    fm := meta.Get(ctx)

    title, ok := fm["title"].(string)
    if !ok || title == "" {
        return Page{}, fmt.Errorf("missing required frontmatter field: title")
    }

    slug, ok := fm["slug"].(string)
    if !ok || slug == "" {
        return Page{}, fmt.Errorf("missing required frontmatter field: slug")
    }

    page := Page{
        Title: title,
        Slug:  slug,
    }

    if subtitle, ok := fm["subtitle"].(string); ok {
        page.Subtitle = subtitle
    }

    page.TOC = buildTOC(doc, src)

    var buf bytes.Buffer
    if err := md.Renderer().Render(&buf, src, doc); err != nil {
        return Page{}, err
    }
    page.Content = template.HTML(buf.String())

    return page, nil
}

type tocItem struct {
    Level int
    Text  string
    ID    string
}

func buildTOC(doc ast.Node, src []byte) template.HTML {
    var items []tocItem

    ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
        if h, ok := n.(*ast.Heading); ok && entering {
            idVal, exists := h.AttributeString("id")
            if !exists {
                return ast.WalkContinue, nil
            }
            id := string(idVal.([]byte))
            text := headingText(h, src)
            items = append(items, tocItem{Level: h.Level, Text: text, ID: id})
        }
        return ast.WalkContinue, nil
    })

    if len(items) == 0 {
        return ""
    }

    var b strings.Builder
    b.WriteString("<details class=\"toc\">\n<summary>Contents</summary>\n<ul>\n")
    for _, item := range items {
        indent := ""
        if item.Level == 3 {
            indent = "  "
        }
        fmt.Fprintf(&b, "%s<li><a href=\"#%s\">%s</a></li>\n", indent, item.ID, item.Text)
    }
    b.WriteString("</ul>\n</details>")

    return template.HTML(b.String())
}

func gitDates(filePath string) (created, updated string) {
    out, err := exec.Command("git", "log", "--follow", "--format=%as", "--", filePath).Output()
    if err != nil || len(strings.TrimSpace(string(out))) == 0 {
        return "", ""
    }
    lines := strings.Split(strings.TrimSpace(string(out)), "\n")
    latest := lines[0]
    oldest := lines[len(lines)-1]
    created = oldest
    if latest != oldest {
        updated = latest
    }
    return
}

func headingText(h *ast.Heading, src []byte) string {
    var b strings.Builder
    for c := h.FirstChild(); c != nil; c = c.NextSibling() {
        switch t := c.(type) {
        case *ast.Text:
            b.Write(t.Segment.Value(src))
        case *ast.String:
            b.Write(t.Value)
        }
    }
    return b.String()
}

func writeHTML(tmpl *template.Template, page Page, outPath string) error {
    if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
        return err
    }

    f, err := os.Create(outPath)
    if err != nil {
        return err
    }
    defer f.Close()

    return tmpl.ExecuteTemplate(f, "layout", page)
}

func copyStatic(srcDir, outDir string) error {
    return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        rel, err := filepath.Rel(srcDir, path)
        if err != nil {
            return err
        }

        dst := filepath.Join(outDir, rel)
        if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
            return err
        }

        data, err := os.ReadFile(path)
        if err != nil {
            return err
        }

        return os.WriteFile(dst, data, 0644)
    })
}
