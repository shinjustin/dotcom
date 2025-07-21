# dotcom

Personal website. Static site generator written in Go.

## Pages

| Route | Source |
|---|---|
| `/` | `content/pages/home.md` |
| `/privacy` | `content/pages/privacy.md` |
| `/terms` | `content/pages/terms.md` |

## Usage

Generate the site:

```bash
go run . generate
```

Preview locally:

```bash
go run . serve
```

Custom port:

```bash
ADDR=:3000 go run . serve
```

Then open [http://localhost:8080](http://localhost:8080) or your custom address.

## Project structure

```
main.go
internal/generator/   # site generator
templates/            # HTML layout templates
static/               # CSS and fonts
content/pages/        # Markdown source files
public/               # generated output (git-ignored)
```

## Content

Edit any file in `content/pages/` and re-run `go run . generate` to rebuild.

Each file requires frontmatter:

```markdown
---
title: Page Title
slug: page-slug
description: Optional description.
---
```

## Configuration

Optional environment variables:

| Variable | Default |
|---|---|
| `CONTENT_DIR` | `content` |
| `TEMPLATE_DIR` | `templates` |
| `STATIC_DIR` | `static` |
| `OUTPUT_DIR` | `public` |
| `BASE_URL` | _(empty)_ |
| `ADDR` | `:8080` |
